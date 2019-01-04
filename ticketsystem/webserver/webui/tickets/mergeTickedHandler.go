package tickets

import (
	"de/vorlesung/projekt/IIIDDD/ticketsystem/logging"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/mail"
	"de/vorlesung/projekt/IIIDDD/ticketsystem/webserver/data/ticket"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

/*
	A ticket merge handler.
 */
type TicketMergeHandler struct {
	Logger        logging.Logger
	TicketContext ticket.TicketContext
	MailContext   mail.MailContext
}

/*
	Merge two tickets.
*/
func (t TicketMergeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "post" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		rawFirstTickedId := r.FormValue("firstTicketId")
		rawSecondTickedId := r.FormValue("secondTicketId")

		firstTicketId, err := strconv.Atoi(rawFirstTickedId)
		if err != nil {
			t.Logger.LogError("TicketMergeHandler", err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}
		secondTicketId, err := strconv.Atoi(rawSecondTickedId)
		if err != nil {
			t.Logger.LogError("TicketMergeHandler", err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		firstTicketExists, firstTicket := t.TicketContext.GetTicketById(firstTicketId)
		if !firstTicketExists {
			t.Logger.LogError("TicketMergeHandler", errors.New("ticket does not exist. id: "+rawFirstTickedId))
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		secondTicketExists, secondTicket := t.TicketContext.GetTicketById(secondTicketId)
		if !secondTicketExists {
			t.Logger.LogError("TicketMergeHandler", errors.New("ticket does not exist. id: "+rawSecondTickedId))
			http.Redirect(w, r, "/", http.StatusBadRequest)
			return

		}

		success, err := t.TicketContext.MergeTickets(firstTicketId, secondTicketId)
		if err != nil {
			t.Logger.LogError("TicketMergeHandler", err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}
		if success {
			olderTicket := firstTicketId
			newerTicket := secondTicketId
			if newerTicket < olderTicket {
				olderTicket = secondTicketId
			}

			// Notify creator of first ticket:
			firstMailSubject := mail.BuildTicketMergeNotificationMailSubject(firstTicket.Info().Id, olderTicket)
			firstMailContent := mail.BuildTicketMergeNotificationMailContent(firstTicket.Info().Creator.Mail,
				firstTicketId, secondTicketId)
			err = t.MailContext.CreateNewOutgoingMail(firstTicket.Info().Creator.Mail, firstMailSubject, firstMailContent)
			if err != nil {
				t.Logger.LogError("TicketAppendMessageHandler", err)
				http.Redirect(w, r, "/ticket/"+strconv.Itoa(olderTicket), http.StatusInternalServerError)
				return
			}
			// Notify creator of second ticket:
			secondMailSubject := mail.BuildTicketMergeNotificationMailSubject(secondTicket.Info().Id, olderTicket)
			secondMailContent := mail.BuildTicketMergeNotificationMailContent(secondTicket.Info().Creator.Mail,
				firstTicketId, secondTicketId)
			err = t.MailContext.CreateNewOutgoingMail(secondTicket.Info().Creator.Mail, secondMailSubject, secondMailContent)
			if err != nil {
				t.Logger.LogError("TicketAppendMessageHandler", err)
				http.Redirect(w, r, "/ticket/"+strconv.Itoa(olderTicket), http.StatusInternalServerError)
				return
			}

			// Redirect to the older ticket:
			http.Redirect(w, r, "/ticket/"+strconv.Itoa(olderTicket), http.StatusOK)

			return
		}
	}
}
