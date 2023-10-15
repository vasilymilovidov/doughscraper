use std::path::PathBuf;
use pitch_detector::{
    note::{detect_note_in_range, NoteDetectionResult},
    pitch::HannedFftDetector
};
use anyhow::Result;
use hound;
use std::{env, fs};

const SAMPLE_RATE: f64 = 44100.0;
const MAX_FREQ: f64 = 1046.50; // C6
const MIN_FREQ: f64 = 32.7; // C1

fn detect_note(sample: Vec<f64>) -> Result<NoteDetectionResult> {
    let mut detector = HannedFftDetector::default();
    let note = detect_note_in_range(&sample, &mut detector, SAMPLE_RATE, MIN_FREQ..MAX_FREQ)
        .ok_or(anyhow::anyhow!("Did not get note"))?;
    Ok(note)
}

fn note(path: PathBuf) -> NoteDetectionResult {
    let mut reader = hound::WavReader::open(path).unwrap();
    let mut f64_samples = Vec::new();
    for result in reader.samples::<i32>() {
        let sample = result.unwrap();
        let normalized_sample = (sample << 8) as f64 / (i32::MAX as f64);
        f64_samples.push(normalized_sample);
    }
    let pitch = detect_note(f64_samples).unwrap();
    pitch
}

fn main() {
    if let Some(dir_path) = env::args().nth(1) {
        // Read the directory content.
        for entry in fs::read_dir(dir_path).unwrap() {
            let entry = entry.unwrap();
            let path = entry.path();

            // Process only files (exclude directories).
            if path.is_file() {
                let old_name = path.file_name().ok_or("File name error").unwrap().to_string_lossy().to_string();

                // Skip non-wave files.
                if old_name.ends_with(".wav") || old_name.ends_with(".mp3") {
                    // Detect the pitch.
                    let pitch = note(path.clone());

                    // Generate a new name.
                    let mut new_name = String::new();
                    match pitch.note_name.to_string().as_str() {
                        "A" => new_name = format!("a{}", pitch.octave.to_string()),
                        "A#" => new_name = format!("bb{}", pitch.octave.to_string()),
                        "B" => new_name = format!("b{}", pitch.octave.to_string()),
                        "C" => new_name = format!("c{}", pitch.octave.to_string()),
                        "C#" => new_name = format!("db{}", pitch.octave.to_string()),
                        "D" => new_name = format!("d{}", pitch.octave.to_string()),
                        "D#" => new_name = format!("eb{}", pitch.octave.to_string()),
                        "E" => new_name = format!("e{}", pitch.octave.to_string()),
                        "F" => new_name = format!("f{}", pitch.octave.to_string()),
                        "F#" => new_name = format!("gb{}", pitch.octave.to_string()),
                        "G" => new_name = format!("g{}", pitch.octave.to_string()),
                        "G#" => new_name = format!("ab{}", pitch.octave.to_string()),
                        _ => {}
                    }
                    let mut result = String::new();
                    if old_name.ends_with(".wav") {
                        result = format!("{}-{}.wav", old_name.trim_end_matches(".wav"), new_name);
                    } else if old_name.ends_with(".mp3") {
                        result = format!("{}-{}.mp3", old_name.trim_end_matches(".mp3"), new_name);
                    }

                    // Create a new path for the file.
                    let parent = path.parent().ok_or("Parent error").unwrap();
                    let new_path = parent.join(result);

                    // Rename the file.
                    fs::rename(path, new_path).unwrap();
                }
            }
        }
    } else {
        println!("Failed to get the directory path.");
    }

}

