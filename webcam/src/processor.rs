extern crate opencv;

use opencv::{
    core::{Point, Vector},
    highgui,
    objdetect,
    prelude::*,
    videoio,
    Result,
};
use regex::Regex;
use std::sync::{Arc, Mutex};
use threadpool::ThreadPool;
use webbrowser;
use rfd::{MessageDialog, MessageDialogResult};

pub struct QrCodeProcessor {
    url_regex: Arc<Regex>,
    pool: ThreadPool,
    processed_qr: Arc<Mutex<Vec<String>>>,
}

impl QrCodeProcessor {
    pub fn new(max_threads: usize, url_pattern: &str) -> Self {
        QrCodeProcessor {
            url_regex: Arc::new(Regex::new(url_pattern).unwrap()),
            pool: ThreadPool::new(max_threads),
            processed_qr: Arc::new(Mutex::new(Vec::new())),
        }
    }

    pub fn start_camera_loop(&self) -> Result<()> {
        // Init Video Camera
        // 0 is the default camera
        let mut cam = videoio::VideoCapture::new(0, videoio::CAP_ANY)?;
        if !videoio::VideoCapture::is_opened(&cam)? {
            eprintln!("Unable to open default camera: {:?}", &cam);
            return Err(opencv::Error::new(opencv::core::StsError, "Camera not opened".to_string()));
        }

        println!("Camera opened successfully!");

        // Initialize QR Code detector
        let detector = objdetect::QRCodeDetector::default()?;

        // Create a window to display video feed
        let window = "QR Code Scanner";
        highgui::named_window(window, highgui::WINDOW_AUTOSIZE)?;

        // Frame processing counter
        let mut total = 0;

        loop {
            let mut frame = Mat::default();

            // Capture a frame from the webcam
            cam.read(&mut frame)?;
            if frame.empty() {
                // Skip processing of empty frames
                continue;
            }

            total += 1;

            // Process every 15th frame
            if total % 15 == 0 {
                self.process_frame(&detector, &frame)?;
            }

            // Display the resulting frame in the window
            highgui::imshow(window, &frame)?;

            // Wait for 1ms and check if the ESC key was pressed to exit
            let key = highgui::wait_key(1)?;
            if key == 27 {
                break;
            }
        }

        Ok(())

    }

    fn process_frame(&self, detector: &objdetect::QRCodeDetector, frame: &Mat) -> Result<()> {
        // Detect multiple QR Codes in frame
        let mut points_v: Vector<Point> = Vector::new();
        let mut decoded_info_v = Vector::new();
        let mut straight_qrcode = Mat::default(); // "Optional" 4th argument (in C++ it is optional at least)

        // Detect multiple QR Codes
        let detected = detector.detect_multi(frame, &mut points_v)?;

        if detected {
            // Decode all detected QR Codes
            let decoded = detector.decode_multi(frame, &points_v, &mut decoded_info_v, &mut straight_qrcode)?;

            if decoded {
                let num = decoded_info_v.len();

                if num > 2 {
                    let mut data_list = Vec::new();

                    for data_ref in &decoded_info_v {
                        let data = data_ref.to_string();
                        data_list.push(data);
                    }

                    for data in data_list {
                        self.handle_qr_code(data);
                    }
                }
            }
        }
        Ok(())
    }

    fn handle_qr_code(&self, data: String) {
        let url = self.url_regex.is_match(&data);
        let data_clone = data.clone();
        let processed_clone = Arc::clone(&self.processed_qr);

        self.pool.execute(move || {
            // Check Code has already been processed
            {
                let mut processed = processed_clone.lock().unwrap();
                if processed.contains(&data_clone) {
                    return;
                }

                processed.push(data_clone.clone());
            }
            if url {
                if MessageDialog::new().set_title("URL Detected")
                .set_description(&format!("QR Code Data:\n{}\n\nDo you want to open this link?\n", data_clone))
                .set_buttons(rfd::MessageButtons::YesNo)
                .show() == MessageDialogResult::Yes {
                    if let Err(e) = webbrowser::open(&data_clone) {
                        eprintln!("Failed to open browser: {:?}", e);
                    }
                }
            };
        })
    }
}