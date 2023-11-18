# go-github.com/bahner/go-space

An Erlang node written in go that handles libp2p for SPACE.

This node will handle the plumbing of messaging and pubsub.

Notably it handles the verification and signing of messages,
so that messages passed to SPACE are believed are trusted.
Since this is the only way in and out of messages for the time
being that should be OK. It's not out-of-scope to move this
to SPACE itself.

Using Vault as a backend for secrets, secrets are easily shared,
but actually SPACE shouldn't need to do this for now.

SPACE has an objectID separate from the IPNS CID. This is
a future proof way of designing the iplementation.
This make rekeying of objects possible. Between go and Elixir
the objectID is used (a NanoID, which is superior to UUID.4).

The NanoID is attached to the DID as the fragment of the did.
NB! The did part is not in any way, shape or form ready. But
this needs to be resolved fairly soon. The did:ipid author is
not responding yet. So a did:space is being thought of.
This is probably not a bad idea any way, as that means we
can add state (IPLD) to the did's.

This way things can be developed separately. SPACE trusts
go and Go needs to be made secure. The gist is that anything
over `lo` is safe and can be trusted. If that doesn't suit you,
then get involved.

The upshot should be that SPACE can be developed naiïvely for
now, without having to sign and verify everything. I believe this
a huge relief.
As stated, this can be changed. Store secrets in Vault and Bob's
Your Uncle™.
