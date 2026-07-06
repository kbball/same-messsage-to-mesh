import { createTheme, alpha } from '@mui/material/styles'

export type ColorMode = 'dark' | 'light'

const brand = {
  50: 'hsl(210, 100%, 95%)',
  100: 'hsl(210, 100%, 92%)',
  200: 'hsl(210, 100%, 80%)',
  300: 'hsl(210, 100%, 65%)',
  400: 'hsl(210, 98%,  48%)',
  500: 'hsl(210, 98%,  42%)',
  600: 'hsl(210, 98%,  55%)',
  700: 'hsl(210, 100%, 35%)',
  800: 'hsl(210, 100%, 16%)',
  900: 'hsl(210, 100%, 21%)',
}

const gray = {
  50: 'hsl(220, 35%, 97%)',
  100: 'hsl(220, 30%, 94%)',
  200: 'hsl(220, 20%, 88%)',
  300: 'hsl(220, 20%, 80%)',
  400: 'hsl(220, 20%, 65%)',
  500: 'hsl(220, 20%, 42%)',
  600: 'hsl(220, 20%, 35%)',
  700: 'hsl(220, 20%, 25%)',
  800: 'hsl(220, 30%,  6%)',
  900: 'hsl(220, 35%,  3%)',
}

export function createAppTheme(mode: ColorMode) {
  const dark = mode === 'dark'

  return createTheme({
    palette: {
      mode,
      primary: {
        main: dark ? brand[600] : brand[500],
        light: brand[300],
        dark: brand[700],
        contrastText: dark ? gray[50] : '#fff',
      },
      background: {
        default: dark ? gray[900] : gray[50],
        paper: dark ? gray[800] : '#fff',
      },
      text: {
        primary: dark ? 'hsl(0, 0%, 100%)' : gray[900],
        secondary: dark ? gray[400] : gray[600],
      },
      divider: dark ? alpha(gray[600], 0.3) : alpha(gray[300], 0.8),
      action: {
        hover: dark ? alpha(gray[600], 0.2) : alpha(gray[200], 0.7),
        selected: dark ? alpha(gray[600], 0.3) : alpha(gray[200], 0.9),
      },
    },
    typography: {
      fontFamily: '"Inter", system-ui, sans-serif',
      h6: { fontSize: '1.125rem', fontWeight: 600, lineHeight: 1.4 },
      subtitle1: { fontSize: '0.875rem', fontWeight: 600 },
      subtitle2: { fontSize: '0.8125rem', fontWeight: 600 },
      body1: { fontSize: '0.875rem' },
      body2: { fontSize: '0.8125rem' },
      caption: { fontSize: '0.75rem' },
      button: { fontSize: '0.875rem', fontWeight: 600, textTransform: 'none' },
    },
    shape: { borderRadius: 8 },
    components: {
      MuiPaper: {
        defaultProps: { elevation: 0 },
        styleOverrides: {
          root: {
            backgroundImage: 'none',
            border: `1px solid ${dark ? alpha(gray[600], 0.3) : alpha(gray[300], 0.8)}`,
          },
        },
      },
      MuiAppBar: {
        defaultProps: { elevation: 0 },
        styleOverrides: {
          root: {
            backgroundImage: 'none',
            backgroundColor: dark ? gray[900] : gray[100],
            borderBottom: `1px solid ${dark ? alpha(gray[600], 0.3) : alpha(gray[300], 0.8)}`,
            color: dark ? 'hsl(0, 0%, 100%)' : gray[900],
          },
        },
      },
      MuiButton: {
        defaultProps: { disableElevation: true },
        styleOverrides: {
          root: { borderRadius: 8, padding: '6px 16px', fontWeight: 600 },
          contained: {
            background: dark ? brand[600] : brand[500],
            '&:hover': { background: brand[700] },
          },
        },
      },
      MuiChip: {
        styleOverrides: { root: { borderRadius: 6, fontWeight: 500 } },
      },
      MuiAlert: {
        styleOverrides: { root: { borderRadius: 8 } },
      },
      MuiTabs: {
        styleOverrides: {
          root: {
            borderBottom: `1px solid ${dark ? alpha(gray[600], 0.3) : alpha(gray[300], 0.8)}`,
          },
        },
      },
      MuiTab: {
        styleOverrides: {
          root: {
            fontWeight: 500,
            fontSize: '0.875rem',
            textTransform: 'none',
            minHeight: 40,
            padding: '8px 16px',
            '&.Mui-selected': { fontWeight: 600 },
          },
        },
      },
      MuiTableHead: {
        styleOverrides: {
          root: {
            '& .MuiTableCell-root': {
              backgroundColor: dark ? alpha(gray[700], 0.5) : gray[100],
              fontWeight: 600,
              fontSize: '0.75rem',
              textTransform: 'uppercase',
              letterSpacing: '0.05em',
              color: dark ? gray[300] : gray[600],
            },
          },
        },
      },
      MuiTableRow: {
        styleOverrides: {
          root: {
            '&:hover': { backgroundColor: dark ? alpha(gray[700], 0.3) : alpha(gray[100], 0.8) },
            '&:last-child td': { borderBottom: 0 },
          },
        },
      },
      MuiTableCell: {
        styleOverrides: {
          root: {
            borderBottom: `1px solid ${dark ? alpha(gray[600], 0.2) : alpha(gray[200], 0.9)}`,
            padding: '8px 12px',
          },
        },
      },
      MuiCssBaseline: {
        styleOverrides: {
          body: {
            backgroundColor: dark ? gray[900] : gray[50],
            scrollbarColor: `${dark ? gray[600] : gray[300]} transparent`,
          },
        },
      },
    },
  })
}
