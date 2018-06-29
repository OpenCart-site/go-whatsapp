package whatsapp_connection

import (
	"fmt"
	"git.willing.nrw/WhatsPoll/whatsapp-connection/binary"
	"git.willing.nrw/WhatsPoll/whatsapp-connection/binary/proto"
	"os"
)

type Handler interface {
	HandleError(err error)
}
type TextMessageHandler interface {
	Handler
	HandleTextMessage(message TextMessage)
}
type ImageMessageHandler interface {
	Handler
	HandleImageMessage(message ImageMessage)
}
type VideoMessageHandler interface {
	Handler
	HandleVideoMessage(message VideoMessage)
}

type JsonMessageHandler interface {
	Handler
	HandleJsonMessage(message string)
}

func (wac *conn) AddHandler(handler Handler) {
	wac.handler = append(wac.handler, handler)
}

func (wac *conn) handle(message interface{}) {
	switch m := message.(type) {
	case error:
		for _, h := range wac.handler {
			go h.HandleError(m)
		}
	case string:
		for _, h := range wac.handler {
			x, ok := h.(JsonMessageHandler)
			if !ok {
				continue
			}
			go x.HandleJsonMessage(m)
		}
	case TextMessage:
		for _, h := range wac.handler {
			x, ok := h.(TextMessageHandler)
			if !ok {
				continue
			}
			go x.HandleTextMessage(m)
		}
	case ImageMessage:
		for _, h := range wac.handler {
			x, ok := h.(ImageMessageHandler)
			if !ok {
				continue
			}
			go x.HandleImageMessage(m)
		}
	case VideoMessage:
		for _, h := range wac.handler {
			x, ok := h.(VideoMessageHandler)
			if !ok {
				continue
			}
			go x.HandleVideoMessage(m)
		}
	}
}

func (wac *conn) dispatch(msg interface{}) {
	if msg == nil || len(wac.handler) == 0 {
		return
	}

	switch message := msg.(type) {
	case *binary.Node:
		if message.Description == "action" {
			if con, ok := message.Content.([]interface{}); ok {
				for a := range con {
					if v, ok := con[a].(*proto.WebMessageInfo); ok {
						wac.dispatchProtoMessage(v)
					}
				}
			}
		}
	case error:
		wac.handle(message)
	case string:
		wac.handle(message)
	default:
		fmt.Fprintf(os.Stderr, "unknown type in dipatcher chan: %T", msg)
	}
}

func (wac *conn) dispatchProtoMessage(msg *proto.WebMessageInfo) {
	switch {

	case msg.GetMessage().GetAudioMessage() != nil:
		//dp.handle(getAudioMessage(msg))

	case msg.GetMessage().GetImageMessage() != nil:
		wac.handle(getImageMessage(msg))

	case msg.GetMessage().GetVideoMessage() != nil:
		wac.handle(getVideoMessage(msg))

	case msg.GetMessage().GetConversation() != "":
		wac.handle(getTextMessage(msg))

	default:
		//cannot match message
	}
}
