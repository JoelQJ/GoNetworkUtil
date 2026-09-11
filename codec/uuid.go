package codec

import "github.com/google/uuid"

func (b *ByteBuf) WriteUUID(u uuid.UUID) {
	b.WriteBytes(u[:])
}

func (b *ByteBuf) ReadUUID() uuid.UUID {
	var u uuid.UUID
	b.ReadBytes(u[:])
	return u
}
