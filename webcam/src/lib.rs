mod processor;

pub use processor::QrCodeProcessor;

/// Start the QR Code detection from the camera
#[no_mangle]
pub extern "C" fn start_qr_code_detection() -> i32 {
    let processor = QrCodeProcessor::new(2,r"^(http|https)://");
    match processor.start_camera_loop() {
        Ok(_) => 0, // Success
        Err(_) => -1, // Failure
    }
}