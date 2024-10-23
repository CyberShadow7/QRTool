mod processor;

pub use processor::QrCodeProcessor;

use opencv::Result;

/// Start the QR Code detection from the camera
#[no_mangle]
pub extern "C" fn start_qr_code_detection() -> Result<()> {
    let processor = QrCodeProcessor::new(2,r"^(http|https)://");
    let _ = processor.start_camera_loop();

    Ok(())
}