import { PipeTransform, Injectable, ArgumentMetadata, BadRequestException } from '@nestjs/common';

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
  xAuthorizations = {
    "token:1": "1",
    "token:2": "2",
    "token:3": "3",
  }

  users = [
    new User("1", "Alice", "https://example.com/alice.jpg"),
    new User("2", "Bob", "https://example.com/bob.jpg"),
    new User("3", "Charlie", "https://example.com/charlie.jpg"),
  ]

  members = [
    new Member("unit1", "1", Role.Owner),
    new Member("unit1", "2", Role.Editor),
    new Member("unit2", "2", Role.Owner),
    new Member("unit2", "3", Role.Editor),
    new Member("unit3", "3", Role.Owner),
    new Member("unit3", "1", Role.Reader),
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
    for (let member of this.members) {
      if (member.unitId === unitId && member.userId === userId) {
        return member;
      }
    }
    return null;
  }

  getMembers(unitId: string): Member[] {
    let members = [];
    for (let member of this.members) {
      if (member.unitId === unitId) {
        members.push(member);
      }
    }
    return members;
  }
}

@Injectable()
export class ParseArrayPipe implements PipeTransform {
  transform(value: any, metadata: ArgumentMetadata) {
    if (!value) {
      return [];
    }
    if (Array.isArray(value)) {
      return value;
    }
    return [value];
  }
}