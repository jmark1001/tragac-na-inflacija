from typing import List

from parsers.abstract_parser import AbstractParser
from parsers.petrol_parser import PetrolParser
from parsers.utilities_parser import UtilitiesParser


def find_parser(sentences: List[str]) -> AbstractParser:
    sentences = " ".join(sentences)
    if "BMB-95" in sentences:
        return PetrolParser()
    elif "ЈКП ИНФОСТАН" in sentences:
        return UtilitiesParser()
    else:
        raise ValueError("No suitable parser found.")
