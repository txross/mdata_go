# Summary
This RFC proposes a standard shared data model and transaction processor implementation of the model to enable shared GS1 standard product data on Hyperledger Sawtooth for the purpose of an intercompany consortium, pre-release of the Hyperledger Grid framework. 

If you want to read about the Grid framework:
 - Click [here](https://github.com/hyperledger/grid-rfcs/pull/5/files) for the RFC
 - Click [here](https://www.hyperledger.org/projects/grid) for the product description

# Use Case
The intercompany consortium is working to automate the resolution of sales & trade claims using Hyperledger Sawtooth distributed ledger technology. If successful, there's potential to expand the ledger to enable track and trace for sales & trade. 

To resolve claims, the ledger must be able to **(a)** privately conduct trade transactions between and among parties and **(b)** publicly share master data to support the private transactions. 

The purpose of this repository is to enable function **b** and provide a transaction processor for use by all nodes in the consortium that can create, update, delete, and otherwise maintain shared product master data publicly on chain. The data supplied by this transaction processor will be data queried by private Sabre trade transactions to validate the business logic of the Sabre smart contracts supporting function **a**. 

# Transaction Processor Specification

## Product Entity
A **__product__** is an archetype of an item that is transacted, traded, or referenced in supply chain. 

For the purposes of the consortium, this product will be a GS1 product identified by a GTIN-14 code. Product attributes will be maintained as key-value pairs within the Product struct. The current design specification limits GTIN to the GTIN-14 specification. The check for invalid GTIN therefore will be limited to checking if the provided GTIN is 14 characters in length. 

The attributes of a Product include:
 - UoM - The unit of measure used when conducting trade (i.e. Cases, LBs)

Unit of Measure is not a GS1 standard attribute, but may be a useful attribute to validate trade transactions. For the purpose of this consortium it is not necessary to restrict product attributes to GS1 standards. 

Other attributes may be added to support supply chain functions. Please review the list of gs1:Product specification [here](https://www.gs1.org/voc/Product) for  GS1 standard product attributes. 



## Transactions
Products are managed by submitting transactions to the Master Data Transaction Procesor specified by the code in this repository. The following transactions are supported:

* ProductCreate - Create a Product and store it in state.
* ProductUpdate - Update (replace) the properties of a Product in state.
* ProductDeactivate - Deactivate a product, setting its state to INACTIVE.
* Product Delete - Remove a Product from state. 

## Permissions
Research is required to enable permissions on GS1 standard products created in this processor. If possible, the products should have ownership and only specified agents of the owners should be able to transact against the product. The Hyperledger Grid framework will achieve this using the Pike processor. 

If as a consortium we decide to add permissions to Product maintenance transactions on this processor, we will need to integrate Pike smart permission functions. You can read about the Pike processor and smart permissions [here](https://sawtooth.hyperledger.org/docs/sabre/nightly/master/smart_permissions.html)

# Reference

## State
Products are stored in state at an address prefixed by the first 6 hex characters of the namespace of the master data processor, `fa3781`.

The next 64 characters will be first 64 charactures of the hex-ecoded representation of a hashed GTIN. The full address then becomes:

Prefix|+|Hashed GTIN
---|---|---
hashed namespace first 6 characters | + | hashed gtin first 64 characters 
`fa3781` | + | `c638b29a67d8b4b3784fb84edadc71367b176a28b29e819f508431d28559a4bc`

## Transaction Payload and Execution

This processor relies on the standard Trasnaction and Batch processing defined [in the official Sawtooth Architecture Guide](https://sawtooth.hyperledger.org/docs/core/nightly/1-1/architecture/transactions_and_batches.html) and implements the go sdk processor (github.com/hyperledger/sawtooth-sdk-go/processor).

### ProductCreate
* Inputs:
    - GTIN-14
    - Attributes in the form of key=value pairs
* Outputs
    - State address of stored product

Invalid Transactions occur in the event of:
 - Invalid GTIN (not one of GTIN-14 spec)
 - Invalid attribute payload (not one of key=value pairs)
 - GTIN already exists

### ProductUpdate
* Inputs:
    - GTIN-14
    - Attributes in the form of key=value pairs
* Outputs
    - State address of stored product

Invalid Transactions occur in the event of:
 - Invalid GTIN (not one of GTIN-14 spec)
 - Invalid attribute payload (not one of key=value pairs)
 - GTIN does not exist

If the transaction submits a GTIN with accompanying attributes that already exist, nothing will happen.

### ProductDeactivate
* Inputs:
    - GTIN-14
* Outputs
    - State address of deactivated product

Invalid Transactions occur in the event of:
 - Invalid GTIN (not one of GTIN-14 spec)
 - GTIN does not exist

### ProductDelete
* Inputs:
    - GTIN-14
* Outputs
    - State address of stored current state

Invalid Transactions occur in the event of:
 - Invalid GTIN (not one of GTIN-14 spec)
 - GTIN does not exist
 - GTIN not in INACTIVE state

 # Future Considerations

 ## Using the Pike processor to determine ownership and agency
  Improvements for this process depend on the ability to integrate the Pike transaction processor and provide organizational metadata and ownership access to GTINS. With organizational metadata, the Product master data transactions could be limited to organizations who have ownership or agency to perform transactions on individual GTINs. Additionally, with the Pike processor each consortium member will have its organizational metadata, including GS1 company prefix, stored on-chain. GS1 prefixes could then be used to further validate GTIN-14 codes submitted to the transaction processor, as opposed to simply validating on length of the code. 

  ## Expansion of Product attributes for further supply chain use cases
  If the claims reconciliation PoC proves successful, there is potentional to further the ledger's use case to also apply to product traceability. At such point, it will be appropriate to include GLN location data in the product or organization metadata. 

  At the claims level, claim validation depends on the product traded between companies, but not the instance of the product. As the consortium becomes more sophisticated in its ledger implementation, it may be prudent to trace product at the batch and unit levels and record more granular data at those levels, such as temperature and humidity during product transport. 