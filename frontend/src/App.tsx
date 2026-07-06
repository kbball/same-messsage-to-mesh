import { useEffect, useMemo, useState } from 'react'
import { BrowserRouter, useNavigate, useLocation } from 'react-router-dom'
import {
  ThemeProvider,
  CssBaseline,
  Box,
  Tabs,
  Tab,
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Tooltip,
} from '@mui/material'
import DarkModeIcon from '@mui/icons-material/DarkMode'
import LightModeIcon from '@mui/icons-material/LightMode'
import { createAppTheme, type ColorMode } from './theme'
import AlertsTab from './ui/pages/AlertsTab'
import FiltersTab from './ui/pages/FiltersTab'
import SDRConfigTab from './ui/pages/SDRConfigTab'
import MQTTConfigTab from './ui/pages/MQTTConfigTab'
import ReferenceDataTab from './ui/pages/ReferenceDataTab'

const TABS = [
  { label: 'Alerts', path: '/alerts' },
  { label: 'Filters', path: '/filters' },
  { label: 'SDR Config', path: '/sdr-config' },
  { label: 'MQTT Config', path: '/mqtt-config' },
  { label: 'Reference Data', path: '/reference-data' },
]

function AppInner() {
  const navigate = useNavigate()
  const { pathname } = useLocation()
  const [colorMode, setColorMode] = useState<ColorMode>('dark')
  const theme = useMemo(() => createAppTheme(colorMode), [colorMode])
  const toggleMode = () => setColorMode((m) => (m === 'dark' ? 'light' : 'dark'))

  useEffect(() => {
    if (pathname === '/') navigate('/alerts', { replace: true })
  }, [pathname, navigate])

  const tab = Math.max(
    0,
    TABS.findIndex((t) => t.path === pathname),
  )

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AppBar position="static">
        <Toolbar variant="dense">
          <Box
            component="img"
            src="/logo.png"
            alt="SAME → Mesh"
            sx={{ height: 40, mr: 1.5, objectFit: 'contain' }}
          />
          <Typography variant="h6" sx={{ flexGrow: 1 }}>
            SAME → Mesh
          </Typography>
          <Tooltip title={colorMode === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}>
            <IconButton onClick={toggleMode} size="small" color="inherit">
              {colorMode === 'dark' ? <LightModeIcon /> : <DarkModeIcon />}
            </IconButton>
          </Tooltip>
        </Toolbar>
        <Tabs
          value={tab}
          onChange={(_, v) => navigate(TABS[v as number].path)}
          textColor="inherit"
          indicatorColor="secondary"
          variant="scrollable"
        >
          {TABS.map(({ label }) => (
            <Tab key={label} label={label} />
          ))}
        </Tabs>
      </AppBar>

      <Box sx={{ p: 2 }}>
        {tab === 0 && <AlertsTab />}
        {tab === 1 && <FiltersTab />}
        {tab === 2 && <SDRConfigTab />}
        {tab === 3 && <MQTTConfigTab />}
        {tab === 4 && <ReferenceDataTab />}
      </Box>
    </ThemeProvider>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <AppInner />
    </BrowserRouter>
  )
}
