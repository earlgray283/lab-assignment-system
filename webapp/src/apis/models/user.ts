export interface ApiUser {
  uid: string;
  gpa: number;
  wishLab: string | null;
  confirmedLab: string | null;
  year: number;
  role: string;
}

export interface UpdateUserPayload {
  userID: string;
}
