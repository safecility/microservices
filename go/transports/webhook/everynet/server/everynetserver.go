package server

import (
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/microservices/go/transports/everynet/messages"
	"io"
	"net/http"
	"os"
	"time"
)

const bearerPrefix = "Bearer "

type EverynetServer struct {
	jwtParser *lib.JWTParser
	uplinks   *pubsub.Topic
}

func NewEverynetServer(jwtParser *lib.JWTParser, uplinks *pubsub.Topic) EverynetServer {
	return EverynetServer{jwtParser: jwtParser, uplinks: uplinks}
}

// Start listen at the given port for /vutility messages
func (en *EverynetServer) Start() {
	handler := http.HandlerFunc(en.handleRequest)

	http.Handle("/eastron", handler)
	//if we're allowed hotdrop raw messages with firmware
	http.Handle("/vutility", handler)

	port, e := os.LookupEnv("PORT")
	if !e {
		port = "8092"
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("could not start http: %v", err))
	}
}

func (en *EverynetServer) handleRequest(w http.ResponseWriter, r *http.Request) {

	err := en.handleAuth(r)
	if err != nil {
		log.Err(err).Msg("could not handle request")
		return
	}

	data, err := io.ReadAll(r.Body)

	if err != nil {
		log.Err(err).Msg("no body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var enMessage messages.ENMessage

	err = json.Unmarshal(data, &enMessage)
	if err != nil {
		log.Error().Err(err).Msg("could not parse enMessage")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go func() {
		err := en.handleLoraMessage(enMessage)
		if err != nil {
			log.Error().Err(err).Msg("could not handle message")
		}
	}()

	w.WriteHeader(200)

	return
}

func (en *EverynetServer) handleLoraMessage(enMessage messages.ENMessage) error {
	switch enMessage.Type {
	case messages.InfoType:
		var info messages.InfoParams
		err := json.Unmarshal(enMessage.Params, &info)
		if err != nil {
			log.Err(err).Msg("could not unmarshall info")
			return err
		}
		log.Debug().Str("info", info.Message).
			Str("code", info.Code).Msg("got info")
		break
	case messages.UplinkType:
		var uplink messages.UplinkParams
		err := json.Unmarshal(enMessage.Params, &uplink)
		if err != nil {
			return err
		}

		now := time.Now()

		sm := &stream.SimpleMessage{
			BrokerDevice: stream.BrokerDevice{
				Source:    "everynet",
				DeviceUID: enMessage.Meta.Device,
			},
			Payload: []byte(uplink.Payload),
			Time:    now,
		}
		_, err = stream.PublishToTopic(sm, en.uplinks)
		if err != nil {
			log.Error().Err(err).Msg("could not publish uplink")
			return err
		}
	case messages.DownlinkType:
		log.Debug().Str("data", string(enMessage.Type)).Msg("downlink")
		break
	case messages.DownlinkRequestType:
		var downlinkRequest messages.DownlinkRequestParams
		err := json.Unmarshal(enMessage.Params, &downlinkRequest)
		if err != nil {
			return err
		}
		log.Debug().Int("counter", downlinkRequest.CounterDown).
			Str("request", fmt.Sprintf("%+v", downlinkRequest)).Msg("got downlink request type")
		break
	default:
		log.Warn().Str("infoType", string(enMessage.Type)).Interface("default", enMessage).Msg("everynet unhandled infoType")
	}
	return nil
}

func (en *EverynetServer) handleAuth(r *http.Request) error {
	auth := r.Header.Get("Authorization")
	log.Debug().Interface("header", auth).Msg("auth")

	if auth == "" || len(auth) < len(bearerPrefix) {
		return fmt.Errorf("invalid authorization header")
	}
	token := auth[len(bearerPrefix):]

	claims, err := en.jwtParser.ParseToken(token)
	if err != nil {
		log.Err(err).Msg("could not parse token")
		return err
	}
	// for the moment we're not interested in the claims
	_ = claims
	return nil
}
