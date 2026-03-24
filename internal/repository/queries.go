package repository

const skipNofS = 4 // skipNofS кол-во пропускаемых кадров стека.

const svAFHistoryOfEntryAndExitAddQuery = "exec dbo.svAF_HistoryOfEntryAndExitAdd ?, ?, ?, ?, ?, ?, ?;" // svAFHistoryOfEntryAndExitAddQuery ХП записывает историю входа и выхода пользователя.

const (
	svPerformerAuthQuery         = "exec dbo.sv_PerformerAuth ?, ?;"      // svPerformerAuthQuery ХП проверяет аутентификацию сотрудника.
	svPerformerFindByTabNumQuery = "exec dbo.sv_PerformerFindByTabNum ?;" // svPerformerFindByTabNumQuery ХП ищет сотрудника по табельному номеру.
)
