package natsstream

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/choria-io/aaasvc/auditors"

	"github.com/choria-io/go-protocol/protocol/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWithGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auditors/NATSStream")
}

var _ = Describe("Auditors/NATSStream", func() {
	Describe("Audit", func() {
		It("Should publish valid audit messages", func() {
			auditor := &NatsStream{
				outbox: make(chan interface{}, 1000),
				site:   "GINKGO",
			}

			rpcreq, err := v1.NewRequest("ginkgo", "ginkgo.example.net", "choria=ginkgo", 60, "9b3a0089dbe0d896c1b79bbc12d61212", "mcollective")
			Expect(err).ToNot(HaveOccurred())
			rpcreq.SetMessage("{}")

			j, err := rpcreq.JSON()
			Expect(err).ToNot(HaveOccurred())

			auditor.Audit(auditors.Allow, "choria=allowed", rpcreq)
			auditor.Audit(auditors.Deny, "choria=denied", rpcreq)

			msg := <-auditor.outbox
			notification, ok := msg.(*Notification)
			Expect(ok).To(BeTrue())
			Expect(notification.Protocol).To(Equal("io.choria.signer.v1.signature_audit"))
			Expect(notification.CallerID).To(Equal("choria=allowed"))
			Expect(notification.Action).To(Equal("allow"))
			Expect(notification.Site).To(Equal("GINKGO"))
			Expect(notification.Time).To(BeNumerically(">=", int64(time.Now().UTC().Unix())))
			Expect(notification.Request).To(Equal(json.RawMessage(j)))

			msg = <-auditor.outbox
			notification, ok = msg.(*Notification)
			Expect(ok).To(BeTrue())
			Expect(notification.Protocol).To(Equal("io.choria.signer.v1.signature_audit"))
			Expect(notification.CallerID).To(Equal("choria=denied"))
			Expect(notification.Action).To(Equal("deny"))
			Expect(notification.Site).To(Equal("GINKGO"))
			Expect(notification.Time).To(BeNumerically(">=", int64(time.Now().UTC().Unix())))

			Expect(notification.Request).To(Equal(json.RawMessage(j)))
		})
	})
})
