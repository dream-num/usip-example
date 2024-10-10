import { Injectable } from '@nestjs/common';

export class User {
  constructor(
    public userId: string,
    public name: string,
    public avatar: string,
  ) {}

  normalize(): {[key: string]: any} {
    return {
      userID: this.userId,
      name: this.name,
      avatar: this.avatar,
    };
  }
}

export enum Role {
  Owner = "owner",
  Editor = "editor",
  Reader = "reader",
}

export class Member {
  constructor(
    public unitId: string,
    public userId: string,
    public role: Role,
  ) {}

  normalize(): {[key: string]: any} {
    return {
      unitID: this.unitId,
      userID: this.userId,
      role: this.role,
    };
  }
}

@Injectable()
export class AppService {
  // Mock data
  xAuthorizations = {
    "token:1": "1",
    "token:2": "2",
    "token:3": "3",
    "token:4": "4",
  }

  // Mock data
  users = [
    new User("1", "Alice", this.generateAvatar("Alice")),
    new User("2", "Bob", this.generateAvatar("Bob")),
    new User("3", "Charlie", this.generateAvatar("Charlie")),
    new User("4", "David", this.generateAvatar("David")),
  ]


  verifyCredential(token: string): string {
    return this.xAuthorizations[token];
  }

  getUser(userId: string): User {
    for (let user of this.users) {
      if (user.userId === userId) {
        return user;
      }
    }
    return null;
  }

  getRole(unitId: string, userId: string): Member {
    // TODO: implement this with real data
    const user = this.getUser(userId);
    
    // Mock data
    if(user.name === "Alice") {
      return new Member(unitId, userId, Role.Owner);
    }
    if(user.name === "Bob" || user.name === "David") {
      return new Member(unitId, userId, Role.Editor);
    }
    if(user.name === "Charlie") {
      return new Member(unitId, userId, Role.Reader);
    }
    return null;
  }

  getMembers(unitId: string): Member[] {
    // TODO: implement this with real data

    // Mock data
    let members = [];
    members.push(new Member(unitId, "1", Role.Owner));
    members.push(new Member(unitId, "2", Role.Editor));
    members.push(new Member(unitId, "3", Role.Reader));
    members.push(new Member(unitId, "4", Role.Editor));
    return members;
  }

  generateAvatar(name: string): string {
    const color = `#${((name.charCodeAt(0) * 12345) & 0xffffff).toString(16).padStart(6, '0')}`;
    
    const svg = `
      <svg xmlns="http://www.w3.org/2000/svg" width="50" height="50" viewBox="0 0 50 50">
        <rect width="50" height="50" fill="${color}" />
        <text x="50%" y="50%" dy=".1em" fill="white" font-family="Arial" font-size="20" text-anchor="middle" dominant-baseline="middle">
          ${name.slice(0, 3)}
        </text>
      </svg>
    `;

    const base64 = Buffer.from(svg).toString('base64');
    return `data:image/svg+xml;base64,${base64}`;
  }
}
