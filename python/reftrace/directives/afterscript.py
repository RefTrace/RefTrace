from dataclasses import dataclass
from ..proto import module_pb2
from .base import Directive

@dataclass(frozen=True)
class AfterScriptDirective(Directive):
    """The 'afterScript' directive specifies a script to run after the main process."""
    _value: module_pb2.AfterScriptDirective

    @property
    def script(self) -> str:
        """The script to execute after the main process."""
        return self._value.script