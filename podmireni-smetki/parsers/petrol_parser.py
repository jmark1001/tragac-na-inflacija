import re
from typing import List, Optional

from parsers.abstract_parser import AbstractParser


class PetrolParser(AbstractParser):
    def get_category(self) -> str:
        return "petrol"

    def get_amount(self, sentences: List[str]) -> Optional[float]:
        amount = None
        candidates = ""
        petrol_anchor, start_idx = "BMB-95", None
        tax_anchor, end_idx = "CE: 17,00", None
        for idx, sentence in enumerate(sentences):
            if petrol_anchor in sentence:
                start_idx = idx + 1
            elif tax_anchor in sentence:
                end_idx = idx
            if start_idx and end_idx:
                candidates = " ".join(sentences[start_idx:end_idx])
                break
        if candidates:
            candidates = candidates.split("x")
            liters = candidates[0]
            price_sum = candidates[1].split()
            price = price_sum[0]
            match = re.findall(r"\d+,\d+", price_sum[1])
            if match:
                amount = float(match[0].replace(",", "."))
            elif liters and price:
                amount = float(liters.replace(",", ".")) * float(
                    price.replace(",", ".")
                )
        return amount
