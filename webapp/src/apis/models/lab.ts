export interface LabList {
  labs: Lab[];
}

export interface Lab {
  id: string;
  name: string;
  capacity: number;
  firstChoice: number;
  secondChoice: number;
  thirdChoice: number;
}
