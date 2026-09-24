
# Decisiones

## D1 ¿Quién pone el estado al crear?
**Opciones:** que lo mande quien crea, o que lo ponga el servidor.
**Qué elegimos:** Que lo mande quien crea, pero validando contra la lista de estados válidos ("RESERVADO", "INGRESADO", "CANCELADO").
**Por qué, en nuestro negocio:** Porque los espacios de parqueo para eventos se pueden emitir con anterioridad con el estado "RESERVADO", pero también se pueden registrar directamente en estado "INGRESADO" si el vehículo llega a pagar directo en la puerta del recinto.
**Qué pasaría con la otra opción:** Si el servidor fijara siempre el estado inicial como "RESERVADO", no se podría registrar el ingreso directo de un vehículo que llega a última hora sin realizar previamente un paso adicional de actualización.