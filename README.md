# PASETO-PLAY
A playground to learn about `PASETO`. `PASETO` has 2 types of `purpose`:
- `local` which uses symmetric-key for encryption 
- `public` which uses asymmetric-key for signing

Using 
- `gin` as the web-server
- `aidanwoods.dev/go-paseto` for the PASETO implementation

## WHAT IS PASETO?
PASETO stands for Platform Agnostic Security Tokens. It allows you to create secure, stateless tokens - kind of like JWTs (JSON Web Tokens). Here's how a PASETO token might look like: 

```
v4.local.dhDLjJ2DLBdr0CVDU8E1UBZtf455vTvBCRKGcSdEnXSjXh6BPv5WToC330WRwEqRg1-JdxHkBifpKJz989vhKS-sjpgA3rzCUZKECTJ4GMgBhbogp_KnL0mCXqQ2poKyU5hBlZps2OxgydzYFF8k1AM0HHxDH4E69pffLyp7V4PFwE6EdSQD6NoFR5MbWu5OvLj8ITueqgDMmJfcfo1ILr5IrFHI8c7QJVNHe8EJTGGR7MTZ30HUmYNMpdGFjAAkke_sEBQIsYD7uidiBAM
```

PASETO tokens are made up of 3 (or 4) parts. The `v4` stands for the version of PASETO being used. `local` is the PASETOs purpose. It can be either `local` or `public`. Then the 3rd part is the `payload`. It's either encrypted (local) or signed (public). Finally, an optional 4th part is the `footer`. It can be used to store additional unencrypted metadata.

So the structure looks something like this:
`Version.Purpose.Payload(.Footer)`

## PASETO LOCAL
Local paseto uses symmetric-key for encrypting your tokens. This is useful when you have a way to store a shared secret safely. An example could be two services running on the same backend server. Another time to use local paseto is when only one service generates and decrypts the token (if you're doing auth on your backend itself instead of making it a separate service) Only your server needs to encrypt and decrypt the data, and hence a local paseto will suffice. IN the short example API I use here, WE DO NOT NEED TO USE PUBLIC PASETO, local will suffice. Public is just for the demo on how you could use it in a system.

## PASETO PUBLIC
Public paseto uses asymmetric-key for signing your tokens. They're best suited for environments where you can't safely share a secret key. Since public pasetos are signed and not encrypted, we use it only for non-sensitive data (you can just run the token through any standard base64 decoder and get the data). A use case for public pasetos would be when you have an auth server (Google) and a separate resource server (Youtube). Google signs, YouTube verifies. YouTube effectively "trusts" Google's signature but cannot create users itself. 

How do you share the public key for your public paseto? You can either save it in a `.env` for every server that needs the public key (tedious) OR you can publish your key to a public URL and the servers can just download the key from there. 

## PERSISTING THE KEYS
I'm persisting the keys by writing them to different files. Everytime the server runs, it checks whether the key file(s) exist. If they do, great, the server can reuse the keys. Otherwise, it will create new keys and write them to the files for persistance. NOTE that in production, these are usually injected via env variables or a secret manager. 

## WHY PASETO?
PASETO is kindof like an upgraded JWT. It gets rid of the `alg` header in JWTs which is a common target for attackers to forge tokens. Instead, it uses a standard cryptographic algorithm that's fixed for the version of PASETO being used (for v4, it's `Ed25519` for public and `XChaCha20` for local). By declaring the purpose of the token in the header, it also defends against key confusion attacks. So by taking the cryptographic choices out of the programmers hands, PASETO ensures security. 

## WHEN NOT TO USE PASETO
When you want to use OAuth2 or OpenID Connect (OIDC). OIDC Standard mandates the use of JWTs, and adding PASETO will add an extra conversion step (from JWT->PASETO) which adds unnecessary complexity. JWTs also have a much wider support for various languages, whereas PASETO doesn't. If you struggle to find a library that compiles on your specific tech stack, JWT is the pragmatic choice.