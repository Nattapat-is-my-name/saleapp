export interface Customer {
  id: string;
  email: string;
  phone: string;
  firstName: string;
  lastName: string;
  address: string;
  notes: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateCustomerInput {
  email: string;
  phone?: string;
  firstName: string;
  lastName: string;
  address?: string;
  notes?: string;
}

export interface UpdateCustomerInput extends Partial<CreateCustomerInput> {}
