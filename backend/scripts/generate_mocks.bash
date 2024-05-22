#!/bin/zsh

set -e

mockgen -source=domain/activity_field.go -destination=mocks/activity_field.go -package=mocks
mockgen -source=domain/company.go -destination=mocks/company.go -package=mocks
mockgen -source=domain/user.go -destination=mocks/user.go -package=mocks
mockgen -source=domain/skill.go -destination=mocks/skill.go -package=mocks
mockgen -source=domain/auth.go -destination=mocks/auth.go -package=mocks
mockgen -source=domain/user_skill.go -destination=mocks/user_skill.go -package=mocks
mockgen -source=domain/fin_report.go -destination=mocks/fin_report.go -package=mocks
mockgen -source=domain/contact.go -destination=mocks/contact.go -package=mocks
mockgen -source=domain/user_activity_field.go -destination=mocks/user_activity_field.go -package=mocks
