import ed25519
import base64

# Generate an Ed25519 keypair
secret_key, public_key = ed25519.create_keypair()

# Convert keys to base64 strings
private_key_base64 = base64.b64encode(secret_key.to_bytes()).decode()
public_key_base64 = base64.b64encode(public_key.to_bytes()).decode()

# Print the keys in base64 format
print("Private Key (Base64):")
print(private_key_base64)

print("Public Key (Base64):")
print(public_key_base64)
