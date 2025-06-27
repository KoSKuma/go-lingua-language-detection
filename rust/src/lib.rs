use lingua::{Language, LanguageDetectorBuilder};
use std::ffi::{CStr, CString};
use std::os::raw::c_char;

#[no_mangle]
pub extern "C" fn detect_language(text: *const c_char) -> *mut c_char {
    let c_str = unsafe {
        assert!(!text.is_null());
        CStr::from_ptr(text)
    };
    
    let text_str = match c_str.to_str() {
        Ok(s) => s,
        Err(_) => return CString::new("error: invalid utf8").unwrap().into_raw(),
    };

    let languages = vec![
        Language::English,
        Language::Spanish,
        Language::French,
        Language::German,
        Language::Italian,
        Language::Portuguese,
        Language::Russian,
        Language::Japanese,
        Language::Korean,
        Language::Chinese,
        Language::Indonesian,
        Language::Malay,
        Language::Thai,
        Language::Vietnamese,
        Language::Tagalog,
    ];

    let detector = LanguageDetectorBuilder::from_languages(&languages).build();
    
    match detector.detect_language_of(text_str) {
        Some(language) => {
            let language_name = format!("{:?}", language);
            CString::new(language_name).unwrap().into_raw()
        }
        None => CString::new("unknown").unwrap().into_raw(),
    }
}

#[no_mangle]
pub extern "C" fn detect_language_with_confidence(text: *const c_char) -> *mut c_char {
    let c_str = unsafe {
        assert!(!text.is_null());
        CStr::from_ptr(text)
    };
    
    let text_str = match c_str.to_str() {
        Ok(s) => s,
        Err(_) => return CString::new("error: invalid utf8").unwrap().into_raw(),
    };

    let languages = vec![
        Language::English,
        Language::Spanish,
        Language::French,
        Language::German,
        Language::Italian,
        Language::Portuguese,
        Language::Russian,
        Language::Japanese,
        Language::Korean,
        Language::Chinese,
        Language::Indonesian,
        Language::Malay,
        Language::Thai,
        Language::Vietnamese,
        Language::Tagalog,
    ];

    let detector = LanguageDetectorBuilder::from_languages(&languages).build();
    
    let results = detector.compute_language_confidence_values(text_str);
    
    if let Some((language, confidence)) = results.first() {
        let result = format!("{:?}:{:.3}", language, confidence);
        CString::new(result).unwrap().into_raw()
    } else {
        CString::new("unknown:0.000").unwrap().into_raw()
    }
}

#[no_mangle]
pub extern "C" fn free_string(ptr: *mut c_char) {
    unsafe {
        if !ptr.is_null() {
            let _ = CString::from_raw(ptr);
        }
    }
} 
