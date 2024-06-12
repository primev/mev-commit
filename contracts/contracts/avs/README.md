
// Write about how v2 (or future version with more decentralization) will 
// give operators the task of doing the pubkey relaying to the mev-commit chain. 
// That is the off-chain process is replaced by operators, who all look for the 
// valset lists posted to some DA layer (eigenDA?), and then race/attest to post
// this to the mev-commit chain. The operator accounts could be auto funded on our chain. 
// Slashing operators in this scheme would require social intervention as it could
// be pretty clear off chain of malicous actions and/or malicious off-chain validation
// of eigenpod conditions, delegation conditions, etc. 

// TODO: Whitelist is now just operators! Every large org seems to have its own operator.
// Note this can be what "operators do" for now. ie. they have the ability to opt-in their users. 
// But we still allow home stakers to opt-in themselves too. 
// Make it very clear that part 2 of opt-in is neccessary to explicitly communicate to 
// the opter-inner that they must follow the relay connection requirement. Otherwise delegators may be 
// blindly frozen. When opting in as a part of step 2, the sender should be running the validators
// its opting in (st. relay requirement is met).