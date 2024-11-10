from pathlib import Path

from handlers.file_handler import process_file

bills_dir = Path(__file__).parent.parent / "samples" / "bills"


def test_file_handler():
    infostan_sample = bills_dir / "Racun_987524.pdf"
    petrol_sample = bills_dir / "petrol-ba.pdf"

    infostan_data = process_file(infostan_sample)
    petrol_data = process_file(petrol_sample)

    assert infostan_data["category"] == "utility", infostan_data["amount"] == 8932.08
    assert petrol_data["category"] == "petrol", petrol_data["amount"] == 141.99
