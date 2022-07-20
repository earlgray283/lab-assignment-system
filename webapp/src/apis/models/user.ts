export interface ApiUser {
  uid: string;
  gpa: number;
  lab1?: string;
  lab2?: string;
  lab3?: string;
  confirmedLab?: string;
}

export interface UserLab {
  lab1: string;
  lab2: string;
  lab3: string;
}
