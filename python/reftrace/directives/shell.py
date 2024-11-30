from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class ShellDirective(Directive):
    """The 'shell' directive specifies the shell to use for script execution."""
    _value: module_pb2.ShellDirective

    @property
    def shell(self) -> str:
        """The shell to use (e.g., 'bash', 'zsh')."""
        return self._value.shell