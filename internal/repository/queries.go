package repository

const svAFHistoryOfEntryAndExitAddQuery = "exec dbo.svAF_HistoryOfEntryAndExitAdd ?, ?, ?, ?, ?, ?, ?;" // svAFHistoryOfEntryAndExitAddQuery ХП записывает историю входа и выхода пользователя.

const (
	svPerformerAuthQuery         = "exec dbo.sv_PerformerAuthWithData ?, ?;"    // svPerformerAuthQuery ХП проверяет аутентификацию сотрудника и хранит данные.
	svPerformerFindByTabNumQuery = "exec dbo.sv_PerformerFindByTabNum ?;"       // svPerformerFindByTabNumQuery ХП ищет сотрудника по табельному номеру.
	svPerformerAllQuery          = "exec dbo.sv_PerformersAll;"                 // svPerformerAllQuery ХП отображает список сотрудников.
	svPerformerUpdQuery          = "exec dbo.sv_PerformerUpd ?, ?, ?, ?, ?, ?;" // svPerformerUpdQuery ХП обновляет сотрудника по ид.
	svPerformerFindByIdQuery     = "exec dbo.sv_PerformerFindById ?;"           // svPerformerFindByIdQuery ХП ищет сотрудника по ид.
)

const svRolesQuery = "exec dbo.sv_Roles;" // svRolesQuery ХП получает список ролей.
