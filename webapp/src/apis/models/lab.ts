export interface LabList {
  labs: Lab[];
}

export interface Lab {
  id: string;
  name: string;
  capacity: number;
  year: number;
  userGPAs: UserGPA[];
}

export interface UserGPA {
  userID: string;
  gpa: number;
}
