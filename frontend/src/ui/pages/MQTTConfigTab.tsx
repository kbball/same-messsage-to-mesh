import { Box, Typography, Paper, Alert, Stack } from '@mui/material'

export default function MQTTConfigTab() {
  return (
    <Box>
      <Typography variant="h6" sx={{ mb: 2 }}>
        MQTT Configuration
      </Typography>
      <Alert severity="info" sx={{ mb: 3 }}>
        MQTT publishing is Phase 2. Configure broker settings via environment variables for now.
      </Alert>
      <Paper sx={{ p: 3, maxWidth: 480 }}>
        <Stack spacing={2}>
          <Typography variant="subtitle1">Environment Variables</Typography>
          <Typography variant="body2" component="div">
            <code>MQTT_ENABLED</code> — set to <code>true</code> to enable publishing
          </Typography>
          <Typography variant="body2" component="div">
            <code>MQTT_HOST</code> — broker hostname (default: localhost)
          </Typography>
          <Typography variant="body2" component="div">
            <code>MQTT_PORT</code> — broker port (default: 1883)
          </Typography>
          <Typography variant="body2" component="div">
            <code>MQTT_PUBLISH_TOPIC</code> — topic to publish alerts to (default: same/alerts)
          </Typography>
        </Stack>
      </Paper>
    </Box>
  )
}
