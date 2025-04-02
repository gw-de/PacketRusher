/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 * © Copyright 2024 Valentin D'Emmanuele
 */
package pdu_session_management

import (
	"encoding/binary"
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"net/netip"

	"github.com/free5gc/aper"
	"github.com/free5gc/ngap"

	customNgapType "my5G-RANTester/lib/ngap/ngapType"

	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
)

type PDUSessionResourceSetupResponseBuilder struct {
	pdu ngapType.NGAPPDU
	ies *ngapType.ProtocolIEContainerPDUSessionResourceSetupResponseIEs
}

func PDUSessionResourceSetupResponse(pduSessions []*context.GnbPDUSession, ue *context.GNBUe, gnb *context.GNBContext) ([]byte, error) {
	return NewPDUSessionResourceSetupResponseBuilder().
		SetAmfUeNgapId(ue.GetAmfUeId()).SetRanUeNgapId(ue.GetRanUeId()).
		SetPDUSessionResourceSetupListSURes(gnb, pduSessions).
		Build()
}

func NewPDUSessionResourceSetupResponseBuilder() *PDUSessionResourceSetupResponseBuilder {
	pdu := ngapType.NGAPPDU{}

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceSetupResponse
	successfulOutcome.Value.PDUSessionResourceSetupResponse = new(ngapType.PDUSessionResourceSetupResponse)

	pDUSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceSetupResponse
	ies := &pDUSessionResourceSetupResponse.ProtocolIEs

	return &PDUSessionResourceSetupResponseBuilder{pdu, ies}
}

func (builder *PDUSessionResourceSetupResponseBuilder) SetAmfUeNgapId(amfUeNgapID int64) *PDUSessionResourceSetupResponseBuilder {
	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	builder.ies.List = append(builder.ies.List, ie)

	return builder
}

func (builder *PDUSessionResourceSetupResponseBuilder) SetRanUeNgapId(ranUeNgapID int64) *PDUSessionResourceSetupResponseBuilder {
	// RAN UE NGAP ID
	ie := ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	builder.ies.List = append(builder.ies.List, ie)

	return builder
}

func (builder *PDUSessionResourceSetupResponseBuilder) SetPDUSessionResourceSetupListSURes(gnb *context.GNBContext, pduSessions []*context.GnbPDUSession) *PDUSessionResourceSetupResponseBuilder {
	// PDU Session Resource Setup Response List
	ie := ngapType.PDUSessionResourceSetupResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceSetupListSURes
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceSetupListSURes
	ie.Value.PDUSessionResourceSetupListSURes = new(ngapType.PDUSessionResourceSetupListSURes)

	pDUSessionResourceSetupListSURes := ie.Value.PDUSessionResourceSetupListSURes

	for _, pduSession := range pduSessions {
		if pduSession == nil {
			continue
		}

		// PDU Session Resource Setup Response Item in PDU Session Resource Setup Response List
		pDUSessionResourceSetupItemSURes := ngapType.PDUSessionResourceSetupItemSURes{}

		// PDU Session ID : This is an unique identifier generated by UE. Can’t be same as any existing PDU session.
		pDUSessionResourceSetupItemSURes.PDUSessionID.Value = pduSession.GetPduSessionId()

		pDUSessionResourceSetupItemSURes.PDUSessionResourceSetupResponseTransfer = GetPDUSessionResourceSetupResponseTransfer(gnb.GetN3GnbIp(), pduSession.GetTeidDownlink(), pduSession.GetQosId())

		pDUSessionResourceSetupListSURes.List = append(pDUSessionResourceSetupListSURes.List, pDUSessionResourceSetupItemSURes)
	}
	builder.ies.List = append(builder.ies.List, ie)

	return builder
}

func (builder *PDUSessionResourceSetupResponseBuilder) SetPDUSessionResourceFailedToSetupListSURes(gnb *context.GNBContext, pduSessions []*context.GnbPDUSession) *PDUSessionResourceSetupResponseBuilder {

	// PDU Sessuin Resource Failed to Setup List
	// ie = ngapType.PDUSessionResourceSetupResponseIEs{}
	// ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceFailedToSetupListSURes
	// ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	// ie.Value.Present = ngapType.PDUSessionResourceSetupResponseIEsPresentPDUSessionResourceFailedToSetupListSURes
	// ie.Value.PDUSessionResourceFailedToSetupListSURes = new(ngapType.PDUSessionResourceFailedToSetupListSURes)

	// pDUSessionResourceFailedToSetupListSURes := ie.Value.PDUSessionResourceFailedToSetupListSURes

	// // PDU Session Resource Failed to Setup Item in PDU Sessuin Resource Failed to Setup List
	// pDUSessionResourceFailedToSetupItemSURes := ngapType.PDUSessionResourceFailedToSetupItemSURes{}
	// pDUSessionResourceFailedToSetupItemSURes.PDUSessionID.Value = 10
	// pDUSessionResourceFailedToSetupItemSURes.PDUSessionResourceSetupUnsuccessfulTransfer = GetPDUSessionResourceSetupUnsucessfulTransfer()

	// pDUSessionResourceFailedToSetupListSURes.List = append(pDUSessionResourceFailedToSetupListSURes.List, pDUSessionResourceFailedToSetupItemSURes)

	// pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)
	// Criticality Diagnostics (optional)
	return builder
}

func (builder *PDUSessionResourceSetupResponseBuilder) Build() ([]byte, error) {
	return ngap.Encoder(builder.pdu)
}

func GetPDUSessionResourceSetupResponseTransfer(ipv4 netip.Addr, teid uint32, qosId int64) []byte {
	data := buildPDUSessionResourceSetupResponseTransfer(ipv4, teid, qosId)
	encodeData, _ := aper.MarshalWithParams(data, "valueExt")
	return encodeData
}

func buildPDUSessionResourceSetupResponseTransfer(ipv4 netip.Addr, teid uint32, qosId int64) (data customNgapType.PDUSessionResourceSetupResponseTransfer) {

	// QoS Flow per TNL Information
	qosFlowPerTNLInformation := &data.QosFlowPerTNLInformation
	qosFlowPerTNLInformation.UPTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel

	// UP Transport Layer Information in QoS Flow per TNL Information
	upTransportLayerInformation := &qosFlowPerTNLInformation.UPTransportLayerInformation
	upTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
	upTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)

	dowlinkTeid := binary.BigEndian.AppendUint32(nil, teid)
	upTransportLayerInformation.GTPTunnel.GTPTEID.Value = dowlinkTeid
	upTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap(ipv4.String(), "")

	// Associated QoS Flow List in QoS Flow per TNL Information
	associatedQosFlowList := &qosFlowPerTNLInformation.AssociatedQosFlowList

	associatedQosFlowItem := ngapType.AssociatedQosFlowItem{}
	associatedQosFlowItem.QosFlowIdentifier.Value = qosId
	associatedQosFlowList.List = append(associatedQosFlowList.List, associatedQosFlowItem)

	return
}
