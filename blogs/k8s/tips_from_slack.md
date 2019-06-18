### LabelSelector
- Don't use `.String()`
- You can use `metav1.FormatLabelSelector(selector)`

### Serializers & Decoders
- Serializers read versioned objects
- Decoders wrap serializers 
- Decoders perform domain specific actions like defaulting and/or converting to internal version
