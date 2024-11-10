from typing import List, Optional

from parsers.abstract_parser import AbstractParser


class UtilitiesParser(AbstractParser):
    def get_category(self) -> str:
        return "utility"

    def get_amount(self, sentences: List[str]) -> Optional[float]:
        amount = None
        amount_anchor = "За уплату: РСД"
        for sentence in sentences:
            if amount_anchor in sentence:
                amount = sentence.replace(amount_anchor, "").replace(".", "").strip()
                break
        return float(amount.replace(",", ".")) if amount else None
