package srv

import (
    "bitbucket.org/rolevax/sakilogy-server/model"
)

type Mail struct {
    To      model.Uid
    Msg     interface{}
}

