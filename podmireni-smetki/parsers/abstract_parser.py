from abc import ABC, abstractmethod
from typing import List, Dict, Optional


class AbstractParser(ABC):

    def parse_data(self, sentences: List[str]) -> Dict:
        return {"category": self.get_category(), "amount": self.get_amount(sentences)}

    @abstractmethod
    def get_category(self) -> str:
        raise NotImplementedError

    @abstractmethod
    def get_amount(self, sentences: List[str]) -> Optional[float]:
        raise NotImplementedError
