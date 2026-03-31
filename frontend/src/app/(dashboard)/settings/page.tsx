"use client";

import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";

export default function SettingsPage() {
  return (
    <Box>
      <Box sx={{ mb: 3 }}>
        <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
          Settings
        </Typography>
        <Typography sx={{ color: "#64748" }}>
          Manage application settings
        </Typography>
      </Box>
      <Card>
        <CardContent>
          <Typography>Settings coming soon...</Typography>
        </CardContent>
      </Card>
    </Box>
  );
}
