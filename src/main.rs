use anyhow::{Context, Result};
use futures_util::StreamExt;
use futures_util::stream::Stream;
use std::sync::Arc;
use tokio::io::AsyncWriteExt;
use tokio::net::TcpStream;
use tokio_rustls::TlsConnector;
use tokio_rustls::rustls;
use tokio_util::codec::{FramedRead, LinesCodec, LinesCodecError};

fn load_root_cert_store_from_file(path: &str) -> Result<rustls::RootCertStore> {
    let mut root_store = rustls::RootCertStore::empty();
    let f = std::fs::File::open(path)
        .context(format!("Failed to open certificate file: {}", path))?;
    let mut reader = std::io::BufReader::new(f);

    // Parse PEM certificates (DER blobs)
    let certs = rustls_pemfile::certs(&mut reader)
        .context("Failed to parse PEM certificates from server.crt")?;
    if certs.is_empty() {
        anyhow::bail!("No PEM certificates found in {}", path);
    }

    // add_parsable_certificates returns (usize, usize) in rustls 0.20: (added, ignored)
    let (added, _ignored) = root_store.add_parsable_certificates(&certs);
    if added == 0 {
        anyhow::bail!("No certificates were added to RootCertStore from {}", path);
    }
    Ok(root_store)
}

#[tokio::main]
async fn main() -> Result<()> {
    // Endereço do servidor TLS (seu servidor Go deve estar rodando nessa porta)
    let addr = "127.0.0.1:2525";
    let domain = "localhost"; // deve bater com o CN do certificado (server.crt)

    // Carrega certificado do servidor (self-signed) para validação local
    // Ajuste o caminho se seu server.crt estiver em outro diretório
    let root_store = load_root_cert_store_from_file("server.crt")
        .context("Loading server.crt into RootCertStore failed")?;

    // Conecta TCP
    let tcp = TcpStream::connect(addr)
        .await
        .context("Failed to connect TCP to server")?;

    // Configura ClientConfig com o root store que contém server.crt
    let config = rustls::ClientConfig::builder()
        .with_safe_defaults()
        .with_root_certificates(root_store)
        .with_no_client_auth();

    let connector = TlsConnector::from(Arc::new(config));

    // Converte nome do host para ServerName (necessário para SNI/validação)
    let server_name = rustls::ServerName::try_from(domain)
        .context("Invalid DNS name for TLS")?;

    // Faz handshake TLS (async)
    let tls_stream = connector
        .connect(server_name, tcp)
        .await
        .context("TLS handshake failed")?;

    // Divide leitura e escrita
    let (r, mut w) = tokio::io::split(tls_stream);

    // FramedRead + LinesCodec para leitura linha-a-linha
    let mut lines = FramedRead::new(r, LinesCodec::new());

    // Lê saudação inicial do servidor (se houver)
    if let Some(line) = lines.next().await {
        let line = line.context("Failed to read initial server response")?;
        println!("Server: {}", line);
    }

    // Exemplos de comandos do protocolo
    send_cmd(&mut w, &mut lines, "EHLO rust-async").await?;
    send_cmd(&mut w, &mut lines, "AUTH user pass").await?;
    send_cmd(&mut w, &mut lines, "MAIL FROM:<user@example.com>").await?;
    send_cmd(&mut w, &mut lines, "RCPT TO:<recipient@example.com>").await?;
    send_cmd(&mut w, &mut lines, "DATA").await?;

    // Envia corpo e finalizador "."
    w.write_all(b"Subject: Async Test\r\n\r\nHello async world!\r\n.\r\n")
        .await
        .context("Failed to write email body")?;
    w.flush().await.context("Failed to flush after body")?;

    if let Some(line) = lines.next().await {
        let line = line.context("Failed to read server response after DATA")?;
        println!("Server: {}", line);
    }

    send_cmd(&mut w, &mut lines, "QUIT").await?;

    Ok(())
}

/// send_cmd: escreve comando e aguarda a próxima linha de resposta.
/// - W: writer que implementa AsyncWriteExt + Unpin
/// - S: stream que produz Result<String, LinesCodecError> (linhas recebidas)
async fn send_cmd<W, S>(w: &mut W, lines: &mut S, cmd: &str) -> Result<()>
where
    W: AsyncWriteExt + Unpin,
    S: Stream<Item = Result<String, LinesCodecError>> + Unpin,
{
    println!("Client: {}", cmd);
    w.write_all(format!("{}\r\n", cmd).as_bytes())
        .await
        .context("Failed to write command")?;
    w.flush().await.context("Failed to flush command")?;

    if let Some(resp) = lines.next().await {
        let resp = resp.context("Failed to read server response")?;
        println!("Server: {}", resp);
    }
    Ok(())
}
