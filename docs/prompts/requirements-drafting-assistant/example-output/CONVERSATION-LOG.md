# Requirements Drafting Conversation: Two-Factor Authentication for Admin Users

## Execution Metadata

**Date:** 2025-12-04
**Agent:** requirements-drafting-assistant
**Initial Request:** "We need to add two-factor authentication for admin users. They've been asking for this for security reasons."
**Duration:** ~20 minutes
**Outcome:** Complete requirement ready for ticketing and Gherkin spec drafting

---

## Full Conversation Transcript

**Business Owner:** "We need to add two-factor authentication for admin users. They've been asking for this for security reasons."

---

**Assistant:** "Thanks for bringing this requirement. Two-factor authentication is definitely important for admin accounts. Let me ask some questions to make sure we understand the full scope and can draft a solid requirement.

First, some context questions:

1. **What prompted this specifically?** Was there a security incident, audit recommendation, or just proactive hardening?

2. **Which admin users?** I see in the business.md that we have user statuses (Active, Inactive, Suspended), but I don't see specific role definitions yet. Are we talking about all users with elevated privileges, or a specific admin role?

3. **When should 2FA be required?** Every login, or only for sensitive operations?

4. **What's driving the timeline?** Is this urgent (compliance deadline) or can we plan this properly?"

---

**Business Owner:** "Good questions. It's proactive - our security team recommended it after their quarterly review. We're talking about users with admin privileges - basically anyone who can create/delete users or change system settings. And yes, we'd want 2FA required at every login for these users. Not urgent, but we'd like to have it in the next quarter."

---

**Assistant:** "Perfect, that helps a lot. Now let's talk about the implementation approach:

1. **2FA method:** What method were you thinking? The most common options are:
   - **SMS codes** (easiest for users but less secure)
   - **Authenticator apps** like Google Authenticator, Authy (more secure, no SMS dependency)
   - **Email codes** (fallback option)
   - **Hardware tokens** like YubiKey (most secure but requires physical devices)

2. **Multiple options:** Should admins be able to choose their preferred method, or do you want to standardize on one?

3. **Enrollment:** How should this work? Should existing admins be forced to set up 2FA on their next login, or give them a grace period?"

---

**Business Owner:** "I think authenticator apps are the sweet spot - secure but not requiring hardware. Let's start with that. For simplicity, we can standardize on authenticator apps for now. Maybe add SMS as a backup option later if users struggle with it.

For enrollment, good question - I think we should give existing admins a 30-day grace period to set it up, but after that they can't log in without 2FA. New admins should be required to set it up during their first login."

---

**Assistant:** "Great choices. That's very practical. Now let me explore some edge cases and security considerations:

**Account Recovery Scenarios:**
1. **What if an admin loses their phone** (and thus their authenticator app)? Do we need a recovery mechanism?
2. **Backup codes:** Should we generate one-time backup codes they can save somewhere safe?
3. **Admin lockout:** If an admin is completely locked out, how do they recover access? Contact super-admin? Manual process?

**Technical Constraints:**
Looking at the architecture.md, I see we use JWT tokens that expire after 24 hours.
4. **How should 2FA interact with JWT tokens?** Should 2FA validation create a special flag in the JWT? Or should 2FA only be required at initial login, with the JWT validity serving as the session?

**Enrollment Edge Cases:**
5. **What happens during the 30-day grace period** if an admin just... never sets it up? Do we send reminders? Block them on day 31?
6. **Can admins disable 2FA once enabled?** Or is it permanent once activated?"

---

**Business Owner:** "Wow, I hadn't thought about some of these! Let me think through them:

1. Recovery mechanism - yes, absolutely needed. We can't have admins completely locked out.

2. Backup codes sound good. Maybe generate 10 single-use codes when they first set up 2FA?

3. For complete lockout, I think we need a super-admin or system administrator role that can reset 2FA for other admins. That person would also need 2FA though!

4. For JWT tokens - I think 2FA should only be required at initial login. Once they pass 2FA and get a valid JWT, that's good for the 24-hour session. Otherwise we'd be asking for codes constantly which would be annoying.

5. During grace period - let's show a banner warning on every login that counts down the days. At day 21, start sending daily email reminders. On day 31, block login until they set it up.

6. Admins should NOT be able to disable 2FA themselves - that defeats the purpose. Only a super-admin should be able to reset someone's 2FA setup (like if they lose their device)."

---

**Assistant:** "Excellent thinking on all of those. The super-admin approach for recovery is smart. Let me probe a few more scenarios to make sure we're covered:

**Security & Audit:**
1. **Should we log 2FA events?** Things like: 2FA setup, successful auth, failed auth attempts, 2FA resets? This would be useful for security auditing.

2. **Failed attempts:** What happens if someone enters wrong 2FA codes repeatedly? Currently your system has rate limiting for login (5 failed attempts per 15 minutes per BR-002). Should 2FA have similar protection?

3. **QR code security:** When setting up authenticator apps, users scan a QR code. Should that QR code expire after a certain time if not used?

**User Experience:**
4. **Setup flow:** When an admin sets up 2FA for the first time, should we require them to successfully verify a code before saving the setup? This ensures they actually got it working.

5. **Device changes:** If an admin gets a new phone, they'd need to set up the authenticator again. Should we prompt them to download new backup codes when they do this?"

---

**Business Owner:** "You're really thorough, I appreciate that!

1. Yes, definitely log all 2FA events. Our security team will want to audit that.

2. Good catch - yes, same rate limiting. 5 failed 2FA attempts should lock them out for 15 minutes. And probably alert our security team if someone's hitting that limit repeatedly.

3. QR code expiring - that makes sense. Maybe 10 minutes? If they start setup but don't complete it, they'd need to restart.

4. Absolutely require verification before saving. I've seen people mess that up and then lock themselves out.

5. Yes, when they reset their 2FA setup, generate fresh backup codes and invalidate the old ones."

---

**Assistant:** "Perfect. A few more questions about the interaction with existing functionality:

**Integration with Current System:**

Looking at the use cases in business.md, I see UC-002 for User Authentication. This would need to be extended.

1. **API impact:** Our authentication API currently returns a JWT on successful email/password. With 2FA, should we:
   - Return a temporary token after password check, requiring a second API call with the 2FA code to get the real JWT?
   - Or return some kind of 'needs 2FA' status and handle it in one call?

2. **REST vs gRPC:** I see from architecture.md we have a gRPC primary API with REST gateway. Should 2FA work the same way for both, or just for REST (user-facing)?

3. **Regular users:** You mentioned this is only for admin users. Regular users will continue with just email/password, correct?

4. **Role changes:** If a regular user gets promoted to admin, should they be immediately required to set up 2FA, or do they get the 30-day grace period?"

---

**Business Owner:** "Good questions on the API side - I'll be honest, some of this is getting more technical than I usually think about. Let me answer what I can:

1. I think the two-step approach makes more sense? Like, they log in with email/password, then we say 'you need 2FA', and they provide the code in a second step. That seems cleaner from a user perspective too.

2. Yes, 2FA should work the same for both gRPC and REST - we want consistent security regardless of how someone accesses the API.

3. Correct - regular users just use email/password. Only admins need 2FA.

4. Ooh, good edge case. If someone gets promoted to admin, I think immediate setup makes sense. Otherwise they could do admin actions during a 30-day window without 2FA, which feels wrong."

---

**Assistant:** "That makes total sense. The two-step API flow is definitely cleaner. And you're right about immediate setup for promotions - much more secure.

Let me ask about one more category - **documentation and communication:**

1. **User documentation:** Will admins need instructions on how to set up their authenticator app? Should we provide recommended apps (Google Authenticator, Authy, Microsoft Authenticator)?

2. **Communication plan:** How will you notify existing admins about this change? Email announcement before the feature launches?

3. **Support:** Should there be a help link or contact during 2FA setup in case admins have trouble?

And finally, **scope boundaries** - let me confirm what's IN and what's OUT of this initial version:

**In scope:**
- 2FA required for admin users at login
- Authenticator app (TOTP) support
- 10 backup codes generated at setup
- 30-day grace period for existing admins (immediate for new/promoted admins)
- Rate limiting on 2FA attempts
- Super-admin can reset 2FA for locked-out admins
- Full audit logging of 2FA events
- QR code expires after 10 minutes
- Must verify code before completing setup

**Out of scope (future enhancements):**
- SMS backup codes
- Email codes
- Hardware token support
- Remember this device for X days

Does that match your understanding?"

---

**Business Owner:** "Yes, exactly! Though for your questions:

1. Documentation - yes, we should show a short help text during setup with recommendations. Let's recommend Google Authenticator and Authy since those are most common.

2. We'll send an email announcement a week before we enable this, and another reminder 3 days before. Our comms team will handle that.

3. Yes, definitely a help link. We can link to our IT support page or maybe show a 'contact your administrator' message.

And your scope summary is perfect. That captures exactly what we need for v1. The 'remember this device' feature would be nice later, but let's get the core working first."

---

**Assistant:** "Excellent. I think we've covered this thoroughly. Let me check a couple of final security considerations that just occurred to me:

1. **Time synchronization:** Authenticator apps rely on time-based codes. If a server's clock is off by more than a minute or two, codes won't validate. Should we add a note in the technical requirements about NTP configuration?

2. **Secret storage:** The 2FA secret (used to generate codes) needs to be stored securely. I see in architecture.md you use standard database storage. Should these secrets be encrypted at rest separately from regular password hashes?

3. **Backup code usage:** When an admin uses one of their 10 backup codes, should we notify them (since it could indicate their account is compromised)? And should we warn them when they're down to their last 2-3 codes?

4. **Grace period bypass:** Is there any emergency situation where we'd need to bypass 2FA? Like if the 2FA system itself is down? Or do we say 2FA must always work - no backdoors?"

---

**Business Owner:** "More great questions:

1. Time sync - yes, please add that to technical requirements. Our DevOps team needs to know.

2. Encryption - yes, 2FA secrets should be encrypted separately. I think this falls under our general 'encrypt sensitive data' policy, but let's make it explicit in the requirements.

3. Backup code notifications - absolutely! Email them when a backup code is used, and warn when they're down to 3 remaining codes. Good security practice.

4. Emergency bypass - hmm. That's a tough one. I think... no backdoors. If 2FA is down, the whole auth system is down anyway, right? But let me think about this. What if we're locked out and need emergency access? Maybe there should be a manual override that requires two super-admins to approve and is heavily logged? Or is that too complex for v1?"

---

**Assistant:** "That's a really thoughtful question. Emergency access is a classic security vs. availability tradeoff. Here's my take:

For v1, I'd suggest: **No emergency bypass** - keep it simple and secure. Here's why:

- If 2FA is down, it's likely due to server issues, not the 2FA mechanism itself
- You'd still have the backup codes for legitimate lockout scenarios
- Adding an emergency bypass adds complexity and potential security holes
- The super-admin can already reset 2FA for locked-out users, which is the escape hatch

However, let's **document this as an open question for the tech team**. They might have insights about system availability we're not considering. They could also propose a secure implementation if they think it's needed.

Does that sound reasonable? We can flag it as 'tech team to review during implementation'?"

---

**Business Owner:** "Yes, that makes sense. Let's go with no bypass for v1, but flag it for the tech team to review. If they come back with a compelling reason why we need one, we can discuss it then. Better to start secure and add capabilities than start too permissive.

I think we've covered everything! This has been really helpful - you brought up a lot of scenarios I hadn't thought about."

---

**Assistant:** "Great collaboration! I've captured a lot of important details through our discussion. Let me summarize everything into a structured requirement document..."

[Full structured requirement output provided - see REVIEW-OUTPUT.md for complete document with 31 acceptance criteria, edge cases, technical considerations, and open questions]

---

**Business Owner:** "This is perfect! You captured everything we talked about, plus organized it way better than I could have. I especially like that you documented the technical considerations and open questions - that'll save a lot of back-and-forth with the dev team.

The edge cases table is really helpful too. I can already see this being useful for our QA team when they write test cases.

Let's move forward with this. I'll create the ticket and route it to the tech lead and requirements analyst as you suggested."

---

**Assistant:** "Excellent! I'm glad this was helpful. A few final recommendations:

1. **Save this conversation transcript** along with the requirement document - it captures the rationale behind decisions, which is valuable if questions come up later.

2. **Share the 'Open Questions' section** with the tech lead early - especially the emergency bypass question. Better to address that before implementation starts.

3. **When QA writes test scenarios**, make sure they cover all the edge cases in the table. The Gherkin features should map directly to these.

4. **Consider a pilot rollout** - maybe test with a small group of admins first before the full 30-day grace period starts?

Good luck with the implementation! Feel free to come back if new requirements emerge during technical review."

---

**Business Owner:** "Great suggestions. I'll definitely do a pilot with our internal IT admins first - good way to find issues before wider rollout. Thanks for all your help!"

---

## Conversation Analysis

### What Worked Well

1. **Systematic exploration** - Assistant asked questions in logical categories (scope, security, UX, integration)
2. **Context-aware** - Referenced existing business rules (BR-002) and architecture (JWT tokens, gRPC)
3. **Proactive edge case identification** - Brought up scenarios BO hadn't considered (role promotions, time sync, backup code notifications)
4. **Security focus** - Consistently asked about security implications (enumeration, rate limiting, encryption)
5. **Scope management** - Clearly defined in/out of v1 to prevent scope creep
6. **Decision documentation** - Captured WHY decisions were made, not just WHAT
7. **Practical handoff** - Identified open questions for tech team and next steps

### Key Insights Discovered Through Conversation

- **Role promotion edge case** - Immediate 2FA setup needed for newly promoted admins
- **Time synchronization dependency** - NTP configuration critical for TOTP
- **Backup code notifications** - Security monitoring for potential compromise
- **Two-step API flow** - Better UX than single-call with complex state
- **No emergency bypass in v1** - Simpler and more secure, flagged for tech review
- **QR code expiry** - Prevents abandoned setups from remaining in memory

### Metrics

- **Questions asked:** ~25 clarifying questions
- **Edge cases identified:** 12 scenarios documented
- **Acceptance criteria generated:** 31 detailed criteria
- **Open questions for tech team:** 4 critical questions
- **Out-of-scope items:** 5 features deferred to future versions
- **BO "I hadn't thought of that" moments:** 6 (role promotions, QR expiry, backup notifications, time sync, secret encryption, emergency bypass)

### Value Delivered

- **Prevented future issues:** Time sync, role promotions, backup code exhaustion would have caused production bugs
- **Ensured security:** Multiple security considerations identified and documented
- **Clear handoff:** Tech team has all context needed to implement
- **Testable requirements:** QA can write test cases directly from edge case table
- **Decision rationale captured:** Future questions can reference this conversation

This conversation demonstrates the value of conversational requirements drafting over static templates - the back-and-forth exploration uncovered critical edge cases that wouldn't appear in a form-based requirements gathering process.
