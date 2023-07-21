export interface LabList {
  labs: Lab[];
}

export interface Lab {
  id: string;
  name: string;
  capacity: number;
  year: number;
  userGPAs: UserGPA[];
  confirmed: boolean;
}

export interface UserGPA {
  gpa: number;
}
