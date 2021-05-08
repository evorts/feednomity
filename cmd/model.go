package cmd

import (
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/jwe"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/mailer"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/view"
)

type library struct {
	db      database.IManager
	mm      mailer.IMailer
	acl     acl.IManager
	logger  logger.IManager
	config  config.IManager
	session session.IManager
	aes     crypt.ICryptAES
	hash    crypt.ICryptHash
	view    view.ITemplateManager
	jwe     jwe.IManager
}
