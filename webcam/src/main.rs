use webcam::start_qr_code_detection;

fn main() {
    match start_qr_code_detection() {
        0 => println!("QR code detection successful!"),
        -1 => eprintln!("There was an error detecting QR codes!"),
        i32::MIN..=-2_i32 | 1_i32..=i32::MAX => println!("Undefined status code returned!"),
    }
}