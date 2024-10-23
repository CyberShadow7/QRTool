use webcam::start_qr_code_detection;

fn main() {
    match start_qr_code_detection() {
        Ok(_) => println!("QR Code detection finished successfully!"),
        Err(e) => eprintln!("Error during QR Code detection: {:?}", e),
    }
}