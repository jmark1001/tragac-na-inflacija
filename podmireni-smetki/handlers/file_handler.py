from parsers import parser_factory
from pathlib import Path
from typing import List, Optional, Dict

import ocrmypdf
import pdfplumber


def ocr_my_pdf(file_path: Path) -> Path:
    ocred_path = file_path.with_stem(file_path.stem + ".ocr")
    ocrmypdf.ocr(file_path, ocred_path, deskew=True)
    return ocred_path


def read_text(ocred_path: Path) -> Optional[List[str]]:
    with pdfplumber.open(ocred_path) as pdf:
        text = pdf.pages[0].extract_text()
    if text:
        return text.split("\n")
    return None


def process_file(file_path: Path) -> Dict:
    sentences = read_text(file_path)
    if not sentences:
        ocred_file = ocr_my_pdf(file_path)
        sentences = read_text(ocred_file)
    parser = parser_factory.find_parser(sentences)
    return parser.parse_data(sentences)
