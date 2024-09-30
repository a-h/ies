# IES playground

The IES ontology is presented in a PDF format which is hard to read and navigate. This repository is an attempt to make the ontology more accessible and interactive.

The `ies.rdf` file is a conversion of the PDF to RDF.

The `main.go` file can process the IES file contents to provide summary stats, and print trees of relationships.

## Print all

```bash
go run .
```

```
<http://ies.data.gov.uk/ontology/ies4>
ies:ExchangePayload
ies:SecurityLabel
ies:andGroup
ies:assessed
ies:attribute
  ies:allocatedSeatNumber
  ies:associatedPersonName
  ...
  ies:vafNumber
ies:groupDescription
ies:groupName
rdf:type
  ies:aCopyOf
  ...
  ies:visaType
rdfs:Class
  ies:ClassOfClassOfElement
    ies:ClassOfClassOfEntity
      ies:ClassOfEntity
        ies:ClassOfAccount
          ies:ClassOfFinancialAccount
        ies:ClassOfAsset
          ies:ClassOfAmountOfMoney
            ies:Currency
          ies:ClassOfDevice
          ...
```

## Print attributes and their children

```bash
go run . -filter attributes -depth 2
```

```
ies:attribute
  ies:allocatedSeatNumber
  ies:associatedPersonName
  ies:confidence
    ies:hmlConfidence
  ies:contactDetailsOnBooking
  ies:currencyAmount
  ies:dialInNumber
  ies:endsIn
  ies:idAuthenticity
  ies:idEmergencyContactName
  ies:idEmergencyContactTelNo
  ies:idFamilyName
  ies:idGivenNames
  ies:ilrProficiency
  ies:iso8601PeriodRepresentation
  ies:issuerIdentificationNumber
  ies:messageContent
  ies:missionPurpose
  ies:natureOfInterest
  ies:objectContentReference
  ies:quantityDelivered
  ies:quantityOffered
  ies:quantityPurchased
  ies:recurrentPeriodRepresentation
  ies:representationValue
    ies:objectContent
  ies:scheduledArrivalTime
  ies:scheduledDepartureTime
  ies:startsIn
  ies:strengthOfInterest
  ies:uriScheme
  ies:uriSchemeName
  ies:vafNumber
```
