class PyomoCompatibleArray:
    def __init__(self, items: list) -> None:
        self.items = items

    def __len__(self) -> int:
        return len(self.items)

    def __getitem__(self, key: int) -> any:
        return self.items[key - 1]

    def __setitem__(self, key: int, value: any) -> None:
        self.items[key - 1] = value

    def __iter__(self):
        return self.items.__iter__()

    def sum(self) -> int:
        return sum(self.items)
