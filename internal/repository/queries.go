package repository

const svAFHistoryOfEntryAndExitAddQuery = "exec dbo.svAF_HistoryOfEntryAndExitAdd ?, ?, ?, ?, ?, ?, ?;" // svAFHistoryOfEntryAndExitAddQuery ХП записывает историю входа и выхода пользователя.

const (
	svPerformerAuthQuery         = "exec dbo.sv_PerformerAuthWithData ?, ?;" // svPerformerAuthQuery ХП проверяет аутентификацию сотрудника и хранит данные.
	svPerformerFindByTabNumQuery = "exec dbo.sv_PerformerFindByTabNum ?;"    // svPerformerFindByTabNumQuery ХП ищет сотрудника по табельному номеру.
	svPerformerAllQuery          = "exec dbo.sv_PerformersAll;"              // svPerformerAllQuery ХП отображает список сотрудников.
	svPerformerUpdQuery          = "exec dbo.sv_PerformerUpd ?, ?, ?, ?, ?, ?;"
)

const svRolesQuery = "exec dbo.sv_Roles;" // svRolesQuery ХП получает список ролей.
