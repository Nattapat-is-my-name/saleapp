"use client";

import * as React from "react";
import TableMUI from "@mui/material/Table";
import TableBodyMUI from "@mui/material/TableBody";
import TableCellMUI from "@mui/material/TableCell";
import TableContainerMUI from "@mui/material/TableContainer";
import TableHeadMUI from "@mui/material/TableHead";
import TableRowMUI from "@mui/material/TableRow";
import PaperMUI from "@mui/material/Paper";

interface TableProps {
  children: React.ReactNode;
  className?: string;
}

function Table({ children, className }: TableProps) {
  return (
    <TableContainerMUI component={PaperMUI} sx={{ borderRadius: 2, boxShadow: 1 }}>
      <TableMUI className={className} size="small">
        {children}
      </TableMUI>
    </TableContainerMUI>
  );
}

function TableHeader({ children, className }: TableProps) {
  return <TableHeadMUI className={className}>{children}</TableHeadMUI>;
}

function TableBody({ children, className }: TableProps) {
  return <TableBodyMUI className={className}>{children}</TableBodyMUI>;
}

interface TableCellProps {
  children?: React.ReactNode;
  className?: string;
}

function TableCell({ children, className }: TableCellProps) {
  return <TableCellMUI className={className}>{children}</TableCellMUI>;
}

function TableRow({ children, className }: TableProps) {
  return <TableRowMUI className={className}>{children}</TableRowMUI>;
}

export { Table, TableHeader, TableBody, TableCell, TableRow };
