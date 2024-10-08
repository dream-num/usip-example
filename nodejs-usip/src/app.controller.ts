import { Controller, Get, Query, Post, Body } from '@nestjs/common';
import { AppService } from './app.service';
import { Cookies } from './decorator/cookies.decorator';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Get('credential')
  verifyCredential(@Cookies('x-authorization') token: string): any {
    let userId = this.appService.verifyCredential(token);
    if (!userId) {
      return { error: "Invalid token" };
    }

    let user = this.appService.getUser(userId);
    if (!user) {
      return { error: "User not found" };
    }

    return {user: user.normalize()};
  }

  @Post('userinfo')
  batchGetUserInfo(@Body("userIDs") userIds: string[]): any {
    let users = [];
    for (let userId of userIds) {
      let user = this.appService.getUser(userId);
      if (user) {
        users.push(user.normalize());
      }
    }
    return {users: users};
  }

  @Get("role")
  getRole(@Query("unitID") unitId: string, @Query("userID") userId: string): any {
    let member = this.appService.getRole(unitId, userId);
    if (!member) {
      return { error: "Member not found" };
    }
    return member.normalize();
  }

  @Post("collaborators")
  getCollaborators(@Body("unitIDs") unitIds: string[]): any {
    let result = {"collaborators": []};
    for (let unitId of unitIds) {
      let unitCollaborator = {unitID: unitId, subjects: []};
      let members = this.appService.getMembers(unitId)
      for (let member of members) {
        let user = this.appService.getUser(member.userId);
        unitCollaborator.subjects.push({
          role: member.role,
          subject: {
            id: user.userId,
            name: user.name,
            avatar: user.avatar,
            type: "user",
          }
        });
        result.collaborators.push(unitCollaborator);
      }
    }
    return result;
  }
}
