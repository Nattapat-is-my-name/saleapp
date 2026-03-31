"use client";

import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";

export default function ReportsPage() {
  return (
    <Box>
      <Box sx={{ mb: 3 }}>
        <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
          Reports
        </Typography>
        <Typography sx={{ color: "#64748b" }}>
          View business reports and analytics
        </Typography>
      </Box>
      <Card>
        <CardContent>
          <Typography>Reports coming soon...</Typography>
        </CardContent>
      </Card>
    </Box>
  );
}
