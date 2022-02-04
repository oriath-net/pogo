# Current status of data file parsing

* **Releases** is the range of major releases that data files were included
  in. "all" means that the file has been observed in all available releases.

* **Structure** is how much work has been done to determine the structure of
  this data file:

  * ✅  means that the file can be parsed without warnings.

  * ⚠️  means that the file parses with warnings, usually indicating that some
    dynamic data is not being processed correctly.

  * ❌  means that the file has a format defined, but it fails to parse the
    current version of the data file.

  * ❓  means that no format is defined for this file.

* **dat64** indicates whether the file's structure is compatible with dat64
  decoding.

  * ✅  means that the current dat64 file can be parsed without warnings, and
    yields the same results as the legacy dat file.

  * ⚠️  means that the dat64 file parses with warnings, or that parsing yields
    different results from the dat file.

  * ❌  means that the dat64 file fails to parse at all.

  * "n/a" means that no dat64 version of this file exists, probably because it
    was removed before the 2.5.0 release.

* **History** indicates whether the structure is consistent with all release
  versions of the data file.

  * ✅  means that all historical versions can be parsed.

  * ❌  means that some or all historical versions fail to parse.

## A

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| AbyssObjects                     | 3.2-     | ✅ | ✅ | ❌
| AbyssRegions                     | 3.2-     | ✅ | ✅ | ✅
| AbyssTheme                       | 3.2-     | ✅ | ✅ | ✅
| AccountQuestFlags                | 3.5-     | ✅ | ✅ | ✅
| AchievementItemRewards           | 3.7-     | ✅ | ✅ | ❌
| AchievementItems                 | all      | ✅ | ✅ | ❌
| Achievements                     | all      | ✅ | ✅ | ❌
| AchievementSetRewards            | 2.0-     | ✅ | ✅ | ❌
| AchievementSets                  | all      | ✅ | ✅ | ✅
| AchievementSetsDisplay           | 2.0-     | ✅ | ✅ | ✅
| ActiveSkills                     | all      | ✅ | ✅ | ❌
| ActiveSkillTargetTypes           | all      | ✅ | ✅ | ✅
| ActiveSkillType                  | all      | ✅ | ✅ | ❌
| Acts                             | 3.14-    | ✅ | ✅ | ✅
| AddBuffToTargetVarieties         | 3.8-     | ✅ | ✅ | ❌
| AdditionalLifeScaling            | 3.5-     | ✅ | ✅ | ✅
| AdditionalLifeScalingPerLevel    | 3.5-     | ✅ | ✅ | ✅
| AdditionalMonsterPacksFromStats  | 3.8-     | ✅ | ✅ | ✅
| AdditionalMonsterPacksStatMode   | 3.8-     | ✅ | ✅ | ✅
| AdvancedSkillsTutorial           | 3.4-     | ✅ | ✅ | ❌
| AegisVariations                  | 3.13-    | ✅ | ✅ | ✅
| AfflictionBalancePerLevel        | 3.10-    | ✅ | ✅ | ✅
| AfflictionEndgameAreas           | 3.10-3.11 | ❓ | ❓ | ❓
| AfflictionEndgameWaveMods        | 3.10-    | ✅ | ✅ | ✅
| AfflictionFixedMods              | 3.10-    | ✅ | ✅ | ❌
| AfflictionRandomModCategories    | 3.10-    | ✅ | ✅ | ✅
| AfflictionRandomSpawns           | 3.10-3.11 | ❓ | ❓ | ❓
| AfflictionRewardMapMods          | 3.10-    | ✅ | ✅ | ❌
| AfflictionRewardTypes            | 3.10-    | ✅ | ✅ | ✅
| AfflictionRewardTypeVisuals      | 3.10-    | ✅ | ✅ | ✅
| AfflictionSplitDemons            | 3.10-    | ✅ | ✅ | ✅
| AfflictionStartDialogue          | 3.10-    | ✅ | ✅ | ✅
| AlternateBehaviourTypes          | 3.10-    | ✅ | ✅ | ✅
| AlternatePassiveAdditions        | 3.7-     | ✅ | ✅ | ✅
| AlternatePassiveSkills           | 3.7-     | ✅ | ✅ | ❌
| AlternateQualityCurrencyDecayFactors | 3.9-3.14 | ✅ | ✅ | ✅
| AlternateQualityTypes            | 3.9-     | ✅ | ✅ | ❌
| AlternateSkillTargetingBehaviours | 3.9-     | ✅ | ✅ | ❌
| AlternateTreePassiveSizes        | 3.7-     | ✅ | ✅ | ✅
| AlternateTreeVersions            | 3.7-     | ✅ | ✅ | ✅
| AnimatedObjectFlags              | 3.16-    | ✅ | ✅ | ✅
| Animation                        | 3.4-     | ✅ | ✅ | ❌
| ApplyDamageFunctions             | 3.6-     | ✅ | ✅ | ✅
| ArchetypeRewards                 | 3.2-     | ✅ | ✅ | ✅
| Archetypes                       | 3.2-     | ✅ | ✅ | ❌
| ArchitectLifeScalingPerLevel     | 3.5-     | ✅ | ✅ | ✅
| ArchitectMapDifficulty           | 3.3-3.4  | ❓ | ❓ | ❓
| ArchnemesisMetaRewards           | 3.17-    | ❓ | ❓ | ❓
| ArchnemesisModComboAchievements  | 3.17-    | ❓ | ❓ | ❓
| ArchnemesisMods                  | 3.17-    | ❓ | ❓ | ❓
| ArchnemesisModVisuals            | 3.17-    | ❓ | ❓ | ❓
| ArchnemesisRecipes               | 3.17-    | ❓ | ❓ | ❓
| AreaInfluenceDoodads             | 3.8-     | ❌ | ❌ | ❌
| AreaTransitionAnimations         | 3.0-     | ✅ | ✅ | ❌
| AreaTransitionAnimationTypes     | 3.0-     | ✅ | ✅ | ✅
| AreaTransitionInfo               | 3.0-     | ✅ | ✅ | ✅
| AreaType                         | 3.0-     | ✅ | ✅ | ✅
| ArmourClasses                    | all      | ✅ | ✅ | ✅
| ArmourSurfaceTypes               | all      | ✅ | ✅ | ✅
| ArmourTypes                      | all      | ✅ | ✅ | ✅
| Ascendancy                       | 2.1-     | ✅ | ✅ | ✅
| AtlasAwakeningStats              | 3.9-3.16 | ✅ | ✅ | ✅
| AtlasBaseTypeDrops               | 3.9-3.16 | ✅ | ✅ | ✅
| AtlasEntities                    | 3.17-    | ❓ | ❓ | ❓
| AtlasExileBossArenas             | 3.9-     | ❌ | ❌ | ✅
| AtlasExileInfluence              | 3.9-     | ❌ | ❌ | ✅
| AtlasExileInfluenceData          | 3.9-3.16 | ✅ | ✅ | ✅
| AtlasExileInfluenceOutcomes      | 3.9-3.16 | ✅ | ✅ | ✅
| AtlasExileInfluenceOutcomeTypes  | 3.9-3.16 | ✅ | ✅ | ✅
| AtlasExileInfluencePacks         | 3.9-3.9  | ❓ | ❓ | ❓
| AtlasExileInfluenceSets          | 3.9-3.16 | ✅ | ✅ | ✅
| AtlasExileRegionQuestFlags       | 3.9-3.16 | ✅ | ✅ | ✅
| AtlasExiles                      | 3.9-     | ❌ | ❌ | ❌
| AtlasFavouredMapSlots            | 3.17-    | ❓ | ❓ | ❓
| AtlasFog                         | 3.6-     | ✅ | ✅ | ✅
| AtlasInfluenceData               | 3.17-    | ❓ | ❓ | ❓
| AtlasInfluenceOutcomes           | 3.2-     | ❌ | ❌ | ✅
| AtlasInfluenceOutcomeTypes       | 3.17-    | ❓ | ❓ | ❓
| AtlasInfluenceSets               | 3.17-    | ❓ | ❓ | ❓
| AtlasMods                        | 3.9-     | ✅ | ✅ | ✅
| AtlasModTiers                    | 3.9-     | ✅ | ✅ | ✅
| AtlasNode                        | 2.4-     | ❌ | ❌ | ❌
| AtlasNodeDefinition              | 3.5-     | ✅ | ✅ | ❌
| AtlasPassiveSkillTreeGroupType   | 3.17-    | ❓ | ❓ | ❓
| AtlasPositions                   | 3.6-     | ✅ | ✅ | ✅
| AtlasPrimordialAltarChoices      | 3.17-    | ❓ | ❓ | ❓
| AtlasPrimordialAltarChoiceTypes  | 3.17-    | ❓ | ❓ | ❓
| AtlasPrimordialBosses            | 3.17-    | ❓ | ❓ | ❓
| AtlasPrimordialBossInfluence     | 3.17-    | ❓ | ❓ | ❓
| AtlasPrimordialBossOptions       | 3.17-    | ❓ | ❓ | ❓
| AtlasQuadrant                    | 3.5-3.13 | ✅ | ✅ | ✅
| AtlasQuestItems                  | 2.4-3.1  | ❓ | ❓ | ❓
| AtlasRegions                     | 3.9-3.16 | ✅ | ✅ | ❌
| AtlasRegionUpgradeRegions        | 3.13-3.16 | ✅ | ✅ | ✅
| AtlasRegionUpgradesInventoryLayout | 3.13-3.16 | ✅ | ✅ | ✅
| AtlasSector                      | 3.5-3.13 | ✅ | ✅ | ✅
| AtlasSkillGraphs                 | 3.13-3.16 | ✅ | ✅ | ✅
| AtlasUpgradesInventoryLayout     | 3.17-    | ❓ | ❓ | ❓
| Attributes                       | all      | ✅ | ✅ | ✅
| AwardDisplay                     | 3.7-     | ✅ | ✅ | ❌

## B

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| BackendErrors                    | all      | ✅ | ✅ | ✅
| BaseItemTypes                    | all      | ❌ | ❌ | ❌
| BestiaryCapturableMonsters       | 3.2-     | ✅ | ✅ | ❌
| BestiaryEncounters               | 3.2-     | ✅ | ✅ | ❌
| BestiaryFamilies                 | 3.2-     | ✅ | ✅ | ❌
| BestiaryGenus                    | 3.2-     | ✅ | ✅ | ❌
| BestiaryGroups                   | 3.2-     | ✅ | ✅ | ✅
| BestiaryNets                     | 3.2-     | ✅ | ✅ | ❌
| BestiaryRecipeComponent          | 3.2-     | ✅ | ✅ | ❌
| BestiaryRecipeItemCreation       | 3.2-3.11 | ❓ | ❓ | ❓
| BestiaryRecipes                  | 3.2-     | ✅ | ❌ | ❌
| BetrayalChoiceActions            | 3.5-     | ✅ | ✅ | ✅
| BetrayalChoices                  | 3.5-     | ✅ | ✅ | ✅
| BetrayalDialogue                 | 3.5-     | ✅ | ✅ | ✅
| BetrayalDialogueCue              | 3.5-     | ✅ | ✅ | ✅
| BetrayalFlags                    | 3.5-     | ✅ | ✅ | ✅
| BetrayalForts                    | 3.5-     | ✅ | ✅ | ✅
| BetrayalJobs                     | 3.5-     | ⚠️ | ❌ | ⚠️
| BetrayalRanks                    | 3.5-     | ✅ | ✅ | ✅
| BetrayalRelationshipState        | 3.5-     | ✅ | ✅ | ✅
| BetrayalTargetFlags              | 3.5-     | ✅ | ✅ | ✅
| BetrayalTargetJobAchievements    | 3.5-     | ✅ | ✅ | ✅
| BetrayalTargetLifeScalingPerLevel | 3.5-     | ✅ | ✅ | ✅
| BetrayalTargets                  | 3.5-     | ❌ | ❌ | ❌
| BetrayalTraitorRewards           | 3.5-     | ✅ | ✅ | ✅
| BetrayalUpgrades                 | 3.5-     | ✅ | ✅ | ✅
| BetrayalUpgradeSlots             | 3.5-     | ✅ | ✅ | ✅
| BetrayalWallLifeScalingPerLevel  | 3.5-     | ✅ | ✅ | ✅
| BeyondDemons                     | 1.2-     | ✅ | ✅ | ✅
| BindableVirtualKeys              | 1.2-     | ✅ | ✅ | ✅
| BlightBalancePerLevel            | 3.8-     | ✅ | ✅ | ✅
| BlightBossLifeScalingPerLevel    | 3.10-    | ✅ | ✅ | ✅
| BlightChestTypes                 | 3.8-     | ✅ | ✅ | ✅
| BlightCraftingItems              | 3.8-     | ✅ | ✅ | ❌
| BlightCraftingRecipes            | 3.8-     | ✅ | ✅ | ✅
| BlightCraftingResults            | 3.8-     | ✅ | ✅ | ✅
| BlightCraftingTypes              | 3.8-     | ✅ | ✅ | ✅
| BlightCraftingUniques            | 3.8-     | ✅ | ✅ | ✅
| BlightedSporeAuras               | 3.8-     | ✅ | ✅ | ✅
| BlightEncounterTypes             | 3.8-     | ✅ | ✅ | ✅
| BlightEncounterWaves             | 3.8-     | ✅ | ✅ | ✅
| BlightRewardTypes                | 3.8-     | ✅ | ✅ | ✅
| BlightStashTabLayout             | 3.11-    | ✅ | ✅ | ❌
| BlightTopologies                 | 3.8-     | ✅ | ✅ | ✅
| BlightTopologyNodes              | 3.8-     | ⚠️ | ❌ | ❌
| BlightTowerAuras                 | 3.8-     | ✅ | ✅ | ✅
| BlightTowers                     | 3.8-     | ✅ | ✅ | ✅
| BlightTowersPerLevel             | 3.8-     | ✅ | ✅ | ✅
| Bloodlines                       | 1.3-     | ✅ | ✅ | ❌
| BloodTypes                       | all      | ✅ | ✅ | ❌
| BonusMasterExp                   | 3.2-3.4  | ❓ | ❓ | ❓
| BreachBossLifeScalingPerLevel    | 3.8-     | ✅ | ✅ | ✅
| BreachElement                    | 3.13-    | ✅ | ✅ | ✅
| BreachstoneUpgrades              | 3.5-     | ✅ | ✅ | ❌
| BuffCategories                   | all      | ✅ | ✅ | ✅
| BuffDefinitions                  | all      | ✅ | ✅ | ❌
| BuffGroups                       | all      | ✅ | ✅ | ✅
| BuffMergeModes                   | all      | ✅ | ✅ | ✅
| BuffStackUIModes                 | 3.5-     | ✅ | ✅ | ✅
| BuffTemplates                    | 3.14-    | ❌ | ❌ | ❌
| BuffVisualOrbArt                 | 3.12-    | ✅ | ✅ | ✅
| BuffVisualOrbs                   | 3.12-    | ✅ | ✅ | ✅
| BuffVisualOrbTypes               | 1.3-     | ❌ | ❌ | ❌
| BuffVisuals                      | all      | ✅ | ✅ | ❌
| BuffVisualsArtVariations         | 3.15-    | ⚠️ | ⚠️ | ⚠️
| BuffVisualSetEntries             | 3.10-    | ✅ | ✅ | ✅
| BuffVisualSets                   | 3.10-    | ✅ | ✅ | ✅

## C

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| CharacterAudioEvents             | all      | ✅ | ✅ | ❌
| CharacterPanelDescriptionModes   | 1.2-     | ✅ | ✅ | ❌
| CharacterPanelStatContexts       | 1.2-     | ✅ | ✅ | ✅
| CharacterPanelStats              | 1.2-     | ⚠️ | ⚠️ | ⚠️
| CharacterPanelTabs               | 1.2-     | ✅ | ✅ | ✅
| Characters                       | all      | ✅ | ✅ | ❌
| CharacterStartItems              | 1.3-3.11 | ❓ | ❓ | ❓
| CharacterStartQuestState         | 2.2-     | ✅ | ✅ | ❌
| CharacterStartStates             | 1.3-     | ✅ | ✅ | ❌
| CharacterStartStateSet           | 2.2-     | ✅ | ✅ | ✅
| CharacterTextAudio               | 2.0-     | ✅ | ✅ | ✅
| ChatIcons                        | 3.13-    | ✅ | ✅ | ✅
| ChestClusters                    | all      | ✅ | ✅ | ✅
| ChestEffects                     | 2.2-     | ✅ | ✅ | ✅
| ChestItemTemplates               | 3.3-3.11 | ❓ | ❓ | ❓
| Chests                           | all      | ❌ | ❌ | ❌
| ClientLeagueAction               | 3.17-    | ❓ | ❓ | ❓
| ClientStrings                    | all      | ✅ | ✅ | ❌
| ClientUIScreens                  | 3.10-    | ✅ | ✅ | ✅
| CloneShot                        | 3.3-     | ✅ | ❌ | ⚠️
| CloneShotVariations              | 2.4-3.2  | ❓ | ❓ | ❓
| Colours                          | 3.16-    | ✅ | ✅ | ✅
| Commands                         | all      | ✅ | ✅ | ✅
| ComponentArmour                  | 0.11-3.14 | ✅ | ✅ | ✅
| ComponentAttributeRequirements   | all      | ✅ | ✅ | ✅
| ComponentCharges                 | all      | ✅ | ✅ | ❌
| CooldownBypassTypes              | all      | ✅ | ✅ | ✅
| CooldownGroups                   | 2.0-     | ✅ | ✅ | ✅
| CostTypes                        | 3.14-    | ✅ | ✅ | ❌
| CraftingBenchCustomActions       | 1.2-     | ✅ | ✅ | ✅
| CraftingBenchOptions             | 1.2-     | ❌ | ❌ | ❌
| CraftingBenchSortCategories      | 3.16-    | ✅ | ✅ | ✅
| CraftingBenchUnlockCategories    | 3.5-     | ✅ | ✅ | ❌
| CraftingItemClassCategories      | 3.5-     | ✅ | ✅ | ✅
| CurrencyItems                    | all      | ❌ | ❌ | ❌
| CurrencyStashTabLayout           | 2.1-     | ❌ | ❌ | ❌
| CurrencyUseTypes                 | all      | ✅ | ✅ | ✅
| CustomLeagueMods                 | 3.4-3.12 | ✅ | ✅ | ❌
| CustomLeagueProperty             | 3.4-3.4  | ❓ | ❓ | ❓

## D

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| DaemonSpawningData               | 3.14-    | ✅ | ✅ | ✅
| DailyMissions                    | 1.3-3.4  | ❓ | ❓ | ❓
| DailyOverrides                   | 1.3-3.4  | ❓ | ❓ | ❓
| DamageHitEffects                 | 3.11-    | ⚠️ | ⚠️ | ❌
| DamageHitTypes                   | 3.14-    | ❓ | ❓ | ❓
| DamageParticleEffects            | 2.3-     | ✅ | ✅ | ❌
| DamageParticleEffectTypes        | 3.6-     | ✅ | ✅ | ✅
| Dances                           | all      | ✅ | ✅ | ✅
| DaressoPitFights                 | 2.0-     | ✅ | ✅ | ✅
| Default                          | all      | ✅ | ✅ | ✅
| DefaultMonsterStats              | all      | ✅ | ✅ | ❌
| DeliriumStashTabLayout           | 3.11-    | ✅ | ✅ | ❌
| DelveAzuriteShop                 | 3.4-     | ✅ | ✅ | ✅
| DelveBiomes                      | 3.4-     | ✅ | ✅ | ❌
| DelveCatchupDepths               | 3.4-     | ✅ | ✅ | ✅
| DelveCraftingModifierDescriptions | 3.4-     | ✅ | ✅ | ✅
| DelveCraftingModifiers           | 3.4-     | ✅ | ✅ | ❌
| DelveCraftingTags                | 3.4-     | ✅ | ✅ | ✅
| DelveDynamite                    | 3.4-     | ✅ | ✅ | ✅
| DelveFeatureRewards              | 3.4-3.11 | ❓ | ❓ | ❓
| DelveFeatures                    | 3.4-     | ⚠️ | ⚠️ | ❌
| DelveFlares                      | 3.4-     | ✅ | ✅ | ✅
| DelveLevelScaling                | 3.4-     | ✅ | ✅ | ❌
| DelveMonsterSpawners             | 3.4-     | ✅ | ✅ | ❌
| DelveResourcePerLevel            | 3.4-     | ✅ | ✅ | ✅
| DelveRewardTierConstants         | 3.14-    | ✅ | ✅ | ✅
| DelveRooms                       | 3.4-     | ✅ | ✅ | ✅
| DelveStashTabLayout              | 3.8-     | ✅ | ✅ | ❌
| DelveUpgrades                    | 3.4-     | ✅ | ✅ | ✅
| DelveUpgradeType                 | 3.4-     | ✅ | ✅ | ✅
| DescentExiles                    | 1.1-     | ✅ | ✅ | ✅
| DescentRewardChests              | all      | ✅ | ✅ | ✅
| DescentStarterChest              | all      | ✅ | ✅ | ❌
| DexIntMissionGuardMods           | 1.2-3.4  | ❓ | ❓ | ❓
| DexIntMissionGuards              | 1.2-3.4  | ❓ | ❓ | ❓
| DexIntMissions                   | 1.2-3.4  | ❓ | ❓ | ❓
| DexIntMissionTargets             | 1.2-3.4  | ❓ | ❓ | ❓
| DexMissionMods                   | 1.2-3.4  | ❓ | ❓ | ❓
| DexMissionMonsters               | 1.2-3.4  | ❓ | ❓ | ❓
| DexMissions                      | 1.2-3.4  | ❓ | ❓ | ❓
| DexMissionTracking               | 1.2-3.4  | ❓ | ❓ | ❓
| DialogueEvent                    | 3.12-    | ✅ | ✅ | ✅
| Difficulties                     | 0.11-2.6 | ✅ | ✅ | ✅
| Directions                       | 3.7-     | ✅ | ✅ | ✅
| DisplayMinionMonsterType         | 2.3-     | ✅ | ✅ | ✅
| DivinationCardArt                | 2.0-     | ✅ | ✅ | ❌
| DivinationCardStashTabLayout     | 2.1-     | ❌ | ❌ | ❌
| Doors                            | 3.8-     | ✅ | ✅ | ✅
| DropEffects                      | 3.4-     | ✅ | ✅ | ✅
| DropModifiers                    | 3.0-3.9  | ✅ | ✅ | ✅
| DropPool                         | all      | ✅ | ✅ | ✅

## E

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| EclipseMods                      | 2.0-     | ✅ | ✅ | ✅
| EffectDrivenSkill                | 3.6-     | ✅ | ✅ | ❌
| Effectiveness                    | all      | ✅ | ✅ | ✅
| EffectivenessCostConstants       | all      | ✅ | ✅ | ❌
| EinharMissions                   | 3.5-     | ✅ | ✅ | ✅
| EinharPackFallback               | 3.5-     | ✅ | ✅ | ✅
| ElderBossArenas                  | 3.2-     | ✅ | ✅ | ✅
| ElderGuardians                   | 3.17-    | ❓ | ❓ | ❓
| ElderMapBossOverride             | 3.2-     | ✅ | ✅ | ✅
| EndlessLedgeChests               | all      | ✅ | ✅ | ✅
| Environments                     | all      | ✅ | ✅ | ❌
| EnvironmentTransitions           | 3.0-     | ✅ | ✅ | ✅
| Essences                         | 2.4-     | ✅ | ✅ | ❌
| EssenceStashTabLayout            | 2.5-     | ✅ | ✅ | ❌
| EssenceType                      | 2.4-     | ✅ | ✅ | ✅
| EventSeason                      | 1.2-3.13 | ✅ | ✅ | ✅
| EventSeasonRewards               | 1.2-3.13 | ✅ | ✅ | ✅
| EvergreenAchievements            | 3.0-     | ✅ | ✅ | ✅
| EvergreenAchievementTypes        | 3.0-     | ✅ | ✅ | ✅
| ExecuteGEAL                      | 3.4-     | ✅ | ✅ | ❌
| ExpandingPulse                   | 3.6-     | ✅ | ✅ | ❌
| ExpeditionAreas                  | 3.15-    | ✅ | ✅ | ✅
| ExpeditionBalancePerLevel        | 3.15-    | ✅ | ✅ | ✅
| ExpeditionCurrency               | 3.15-    | ✅ | ✅ | ❌
| ExpeditionDealFamilies           | 3.16-    | ❓ | ❓ | ❓
| ExpeditionDeals                  | 3.16-    | ✅ | ✅ | ✅
| ExpeditionFactions               | 3.15-    | ✅ | ✅ | ✅
| ExpeditionMarkersCommon          | 3.15-    | ✅ | ✅ | ✅
| ExpeditionNPCs                   | 3.15-    | ✅ | ✅ | ✅
| ExpeditionRelicModCategories     | 3.15-    | ❓ | ❓ | ❓
| ExpeditionRelicMods              | 3.15-    | ✅ | ✅ | ✅
| ExpeditionRelics                 | 3.15-    | ✅ | ✅ | ✅
| ExpeditionStorageLayout          | 3.15-    | ✅ | ✅ | ❌
| ExpeditionTerrainFeatures        | 3.15-    | ✅ | ✅ | ✅
| ExperienceLevels                 | all      | ✅ | ✅ | ✅
| ExplodingStormBuffs              | 2.0-     | ✅ | ✅ | ❌
| ExtraTerrainFeatureFamily        | 2.4-3.7  | ❓ | ❓ | ❓
| ExtraTerrainFeatures             | 2.3-     | ✅ | ✅ | ❌

## F

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| FixedHideoutDoodads              | 3.5-3.7  | ❓ | ❓ | ❓
| FixedHideoutDoodadTypes          | 3.8-     | ✅ | ✅ | ✅
| FixedMissions                    | 3.5-     | ✅ | ✅ | ✅
| Flasks                           | all      | ✅ | ✅ | ❌
| FlaskStashBaseTypeOrdering       | 3.17-    | ❓ | ❓ | ❓
| FlaskType                        | all      | ✅ | ✅ | ✅
| FlavourText                      | all      | ✅ | ✅ | ✅
| FlavourTextImages                | 3.0-     | ✅ | ✅ | ✅
| Footprints                       | 2.1-     | ❌ | ❌ | ✅
| FootstepAudio                    | 3.8-     | ✅ | ✅ | ✅
| FragmentStashTabLayout           | 3.2-     | ❌ | ❌ | ❌

## G

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| GameConstants                    | all      | ✅ | ✅ | ❌
| GamepadButton                    | 3.16-    | ✅ | ✅ | ✅
| GamepadButtonCombination         | 3.17-    | ❓ | ❓ | ❓
| GamepadThumbstick                | 3.17-    | ❓ | ❓ | ❓
| GamepadType                      | 3.16-    | ✅ | ✅ | ✅
| GameStats                        | 3.8-     | ✅ | ✅ | ✅
| GemTags                          | all      | ✅ | ✅ | ❌
| GemTypes                         | all      | ✅ | ✅ | ✅
| GenericBuffAuras                 | 3.10-    | ✅ | ✅ | ✅
| GenericLeagueRewardTypes         | 3.17-    | ❓ | ❓ | ❓
| GenericLeagueRewardTypeVisuals   | 3.17-    | ❓ | ❓ | ❓
| GeometryAttack                   | 3.3-     | ❌ | ❌ | ❌
| GeometryAttackShapes             | 3.2-3.2  | ❓ | ❓ | ❓
| GeometryAttackTargetTypes        | 3.2-3.2  | ❓ | ❓ | ❓
| GeometryAttackVariations         | 3.2-3.2  | ❓ | ❓ | ❓
| GeometryChannel                  | 3.12-    | ✅ | ✅ | ❌
| GeometryProjectiles              | 3.5-     | ✅ | ✅ | ❌
| GeometryTrigger                  | 3.6-     | ❌ | ❌ | ❌
| GiftWrapArtVariations            | 3.13-    | ✅ | ✅ | ✅
| GlobalAudioConfig                | all      | ✅ | ✅ | ❌
| GlobalQuestFlags                 | 2.3-2.6  | ❓ | ❓ | ❓
| Grandmasters                     | 1.3-     | ✅ | ✅ | ❌
| GrantedEffectGroups              | 2.3-     | ✅ | ✅ | ✅
| GrantedEffectQualityStats        | 3.12-    | ⚠️ | ❌ | ❌
| GrantedEffectQualityTypes        | 3.12-    | ✅ | ✅ | ✅
| GrantedEffects                   | all      | ⚠️ | ❌ | ❌
| GrantedEffectsPerLevel           | all      | ✅ | ✅ | ❌
| GroundEffectEffectTypes          | 3.12-    | ✅ | ✅ | ✅
| GroundEffects                    | 3.4-     | ❌ | ❌ | ❌
| GroundEffectTypes                | 3.4-     | ✅ | ✅ | ❌

## H

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| HarbingerMaps                    | 3.0-3.11 | ❓ | ❓ | ❓
| Harbingers                       | 3.0-     | ✅ | ✅ | ✅
| HarvestColours                   | 3.11-    | ✅ | ✅ | ✅
| HarvestCraftOptionIcons          | 3.12-    | ✅ | ✅ | ✅
| HarvestCraftOptions              | 3.11-    | ✅ | ✅ | ❌
| HarvestCraftTiers                | 3.11-    | ✅ | ✅ | ✅
| HarvestDurability                | 3.11-3.12 | ✅ | ✅ | ✅
| HarvestEncounterScaling          | 3.11-    | ✅ | ✅ | ✅
| HarvestInfrastructure            | 3.11-    | ✅ | ✅ | ❌
| HarvestInfrastructureCategories  | 3.11-3.12 | ✅ | ✅ | ✅
| HarvestMetaCraftingOptions       | 3.11-    | ✅ | ✅ | ✅
| HarvestObjects                   | 3.11-3.12 | ✅ | ✅ | ✅
| HarvestPerLevelValues            | 3.11-    | ✅ | ✅ | ❌
| HarvestPlantBoosterFamilies      | 3.12-3.12 | ✅ | ✅ | n/a
| HarvestPlantBoosters             | 3.12-3.12 | ✅ | ✅ | n/a
| HarvestSeeds                     | 3.13-    | ⚠️ | ⚠️ | ⚠️
| HarvestSeedTypes                 | 3.11-3.12 | ✅ | ✅ | ✅
| HarvestSpecialCraftCosts         | 3.11-3.12 | ✅ | ✅ | ✅
| HarvestSpecialCraftOptions       | 3.11-3.12 | ✅ | ✅ | ✅
| HarvestStorageLayout             | 3.11-3.12 | ✅ | ✅ | ❌
| HeistAreaFormationLayout         | 3.12-    | ✅ | ✅ | ✅
| HeistAreas                       | 3.12-    | ✅ | ✅ | ✅
| HeistBalancePerLevel             | 3.12-    | ✅ | ✅ | ❌
| HeistBlueprintWindowTypes        | 3.12-    | ✅ | ✅ | ✅
| HeistChestRewardTypes            | 3.12-    | ⚠️ | ⚠️ | ❌
| HeistChests                      | 3.12-    | ✅ | ✅ | ❌
| HeistChestTypes                  | 3.12-    | ✅ | ✅ | ✅
| HeistChokepointFormation         | 3.12-    | ✅ | ✅ | ❌
| HeistConstants                   | 3.12-    | ✅ | ✅ | ✅
| HeistContracts                   | 3.12-    | ✅ | ✅ | ✅
| HeistDoodadNPCs                  | 3.12-    | ✅ | ✅ | ✅
| HeistDoors                       | 3.12-    | ✅ | ✅ | ✅
| HeistEquipment                   | 3.12-    | ✅ | ✅ | ✅
| HeistFormationMarkerType         | 3.12-    | ✅ | ✅ | ✅
| HeistGeneration                  | 3.12-    | ✅ | ✅ | ❌
| HeistIntroAreas                  | 3.12-    | ✅ | ✅ | ✅
| HeistJobs                        | 3.12-    | ✅ | ✅ | ❌
| HeistJobsExperiencePerLevel      | 3.12-    | ✅ | ✅ | ✅
| HeistLockType                    | 3.12-    | ✅ | ✅ | ✅
| HeistNPCAuras                    | 3.12-    | ✅ | ✅ | ✅
| HeistNPCBlueprintTypes           | 3.12-    | ✅ | ✅ | ✅
| HeistNPCDialogue                 | 3.12-    | ✅ | ✅ | ✅
| HeistNPCs                        | 3.12-    | ✅ | ✅ | ❌
| HeistNPCStats                    | 3.12-    | ✅ | ✅ | ✅
| HeistObjectives                  | 3.12-    | ✅ | ✅ | ✅
| HeistObjectiveValueDescriptions  | 3.13-    | ✅ | ✅ | ✅
| HeistPatrolPacks                 | 3.12-    | ✅ | ✅ | ✅
| HeistQuestContracts              | 3.12-    | ✅ | ✅ | ❌
| HeistRevealingNPCs               | 3.12-    | ✅ | ✅ | ✅
| HeistRooms                       | 3.12-    | ✅ | ✅ | ✅
| HeistRoomTypes                   | 3.12-    | ✅ | ✅ | ✅
| HeistStorageLayout               | 3.12-    | ✅ | ✅ | ✅
| HeistValueScaling                | 3.12-    | ✅ | ✅ | ✅
| HellscapeAOReplacements          | 3.16-    | ✅ | ✅ | ✅
| HellscapeAreaPacks               | 3.16-    | ✅ | ✅ | ✅
| HellscapeExperienceLevels        | 3.16-    | ✅ | ✅ | ✅
| HellscapeFactions                | 3.16-    | ✅ | ✅ | ✅
| HellscapeImmuneMonsters          | 3.16-    | ✅ | ✅ | ✅
| HellscapeItemModificationTiers   | 3.16-    | ✅ | ✅ | ✅
| HellscapeLifeScalingPerLevel     | 3.16-    | ✅ | ✅ | ✅
| HellscapeModificationInventoryLayout | 3.16-    | ✅ | ✅ | ✅
| HellscapeMods                    | 3.16-    | ❌ | ❌ | ✅
| HellscapeMonsterPacks            | 3.16-    | ✅ | ✅ | ✅
| HellscapePassives                | 3.16-    | ✅ | ✅ | ✅
| HellscapePassiveTree             | 3.16-    | ✅ | ✅ | ✅
| HideoutDoodadCategory            | 3.17-    | ❓ | ❓ | ❓
| HideoutDoodads                   | 1.2-     | ❌ | ❌ | ❌
| HideoutDoodadTags                | 3.16-    | ❌ | ❌ | ✅
| HideoutInteractable              | 1.2-2.6  | ❓ | ❓ | ❓
| HideoutNPCs                      | 3.3-     | ❌ | ❌ | ❌
| HideoutRarity                    | 3.5-     | ✅ | ✅ | ✅
| Hideouts                         | 1.2-     | ✅ | ✅ | ❌

## I

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| ImpactSoundData                  | all      | ✅ | ✅ | ✅
| Incubators                       | 3.7-     | ✅ | ✅ | ✅
| IncursionArchitect               | 3.3-     | ✅ | ✅ | ✅
| IncursionBrackets                | 3.3-     | ✅ | ✅ | ✅
| IncursionChestRewards            | 3.3-     | ✅ | ✅ | ✅
| IncursionChests                  | 3.3-     | ✅ | ✅ | ✅
| IncursionRoomAdditionalBossDrops | 3.3-3.11 | ❓ | ❓ | ❓
| IncursionRoomBossFightEvents     | 3.3-     | ✅ | ✅ | ✅
| IncursionRooms                   | 3.3-     | ⚠️ | ❌ | ❌
| IncursionUniqueUpgradeComponents | 3.3-     | ✅ | ✅ | ✅
| IncursionUniqueUpgrades          | 3.3-3.11 | ❓ | ❓ | ❓
| IndexableSupportGems             | 3.11-    | ✅ | ✅ | ✅
| InfluenceExalts                  | 3.9-     | ✅ | ✅ | ✅
| InfluenceModUpgrades             | 3.13-    | ✅ | ✅ | ❌
| InfluenceTags                    | 3.14-    | ✅ | ✅ | ✅
| InfluenceTypes                   | 3.9-     | ✅ | ✅ | ✅
| IntMissionMods                   | 1.2-3.4  | ❓ | ❓ | ❓
| IntMissionMonsters               | 1.2-3.4  | ❓ | ❓ | ❓
| IntMissions                      | 1.2-3.4  | ❓ | ❓ | ❓
| InvasionMonsterGroups            | 1.1-     | ✅ | ✅ | ✅
| InvasionMonsterRestrictions      | 1.1-     | ⚠️ | ❌ | ❌
| InvasionMonsterRoles             | 1.1-     | ✅ | ✅ | ✅
| InvasionMonstersPerArea          | 1.1-     | ✅ | ✅ | ❌
| Inventories                      | 3.7-     | ✅ | ✅ | ❌
| InventoryId                      | 3.7-3.14 | ✅ | ✅ | ❌
| InventoryType                    | 3.7-     | ✅ | ✅ | ❌
| ItemClassCategories              | 3.0-     | ✅ | ✅ | ❌
| ItemClasses                      | all      | ❌ | ❌ | ❌
| ItemClassesDisplay               | 1.2-2.1  | ✅ | n/a | ✅
| ItemCostPerLevel                 | 3.12-    | ✅ | ✅ | ❌
| ItemCosts                        | 3.16-    | ✅ | ✅ | ✅
| ItemCreationTemplateCustomAction | 3.7-     | ✅ | ✅ | ✅
| ItemExperiencePerLevel           | all      | ✅ | ✅ | ✅
| ItemisedVisualEffect             | all      | ❌ | ❌ | ❌
| ItemNoteCode                     | 3.12-    | ✅ | ✅ | ✅
| ItemSetNames                     | 3.3-     | ✅ | ✅ | ✅
| ItemShopType                     | 3.3-3.16 | ✅ | ✅ | ❌
| ItemStances                      | 3.14-    | ✅ | ✅ | ✅
| ItemSynthesisCorruptedMods       | 3.6-3.6  | ❓ | ❓ | ❓
| ItemSynthesisMods                | 3.6-3.6  | ❓ | ❓ | ❓
| ItemThemes                       | 3.0-3.16 | ✅ | ✅ | ✅
| ItemTradeData                    | 3.3-3.12 | ❓ | ❓ | ❓
| ItemVisualEffect                 | all      | ✅ | ✅ | ❌
| ItemVisualHeldBodyModel          | 3.6-     | ✅ | ✅ | ✅
| ItemVisualIdentity               | all      | ✅ | ✅ | ❌

## J

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| JobAssassinationSpawnerGroups    | 3.5-     | ✅ | ✅ | ✅
| JobRaidBrackets                  | 3.5-     | ✅ | ✅ | ✅

## K

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| KillstreakThresholds             | 1.2-     | ✅ | ✅ | ✅
| KiracLevels                      | 3.17-    | ❓ | ❓ | ❓

## L

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| LabyrinthAreas                   | 2.1-     | ✅ | ✅ | ❌
| LabyrinthBonusItems              | 3.11-    | ✅ | ✅ | ✅
| LabyrinthExclusionGroups         | 2.1-     | ✅ | ✅ | ✅
| LabyrinthIzaroChests             | 2.2-     | ✅ | ✅ | ✅
| LabyrinthLadderRewards           | 2.3-2.6  | ❓ | ❓ | ❓
| LabyrinthNodeOverrides           | 2.2-     | ⚠️ | ❌ | ⚠️
| LabyrinthRewards                 | 2.2-3.11 | ❓ | ❓ | ❓
| LabyrinthRewardTypes             | 3.0-     | ✅ | ✅ | ✅
| Labyrinths                       | 3.0-     | ✅ | ✅ | ❌
| LabyrinthSecretEffects           | 2.1-     | ✅ | ✅ | ❌
| LabyrinthSecretLocations         | 2.1-     | ✅ | ✅ | ✅
| LabyrinthSecrets                 | 2.1-     | ✅ | ⚠️ | ❌
| LabyrinthSection                 | 2.1-     | ✅ | ✅ | ❌
| LabyrinthSectionLayout           | 2.1-     | ✅ | ✅ | ❌
| LabyrinthTrials                  | 3.0-     | ✅ | ✅ | ❌
| LabyrinthTrinkets                | 2.2-     | ✅ | ✅ | ❌
| Languages                        | 2.1-     | ❓ | n/a | ❓
| LeagueCategory                   | all      | ✅ | ✅ | ✅
| LeagueFlag                       | all      | ❌ | ❌ | ❌
| LeagueFlags                      | 3.4-3.14 | ✅ | ✅ | ✅
| LeagueInfo                       | 3.4-     | ❌ | ❌ | ❌
| LeagueQuestFlags                 | 1.2-     | ✅ | ✅ | ✅
| LeagueStaticRewards              | 3.17-    | ❓ | ❓ | ❓
| LeagueTrophy                     | 2.1-3.12 | ✅ | ✅ | ✅
| LegacyAtlasInfluenceOutcomes     | 3.17-    | ❓ | ❓ | ❓
| LegionBalancePerLevel            | 3.7-     | ✅ | ✅ | ✅
| LegionChestCounts                | 3.7-     | ✅ | ✅ | ✅
| LegionChests                     | 3.7-     | ✅ | ✅ | ✅
| LegionFactions                   | 3.7-     | ✅ | ✅ | ❌
| LegionMonsterCounts              | 3.7-     | ✅ | ✅ | ✅
| LegionMonsterTypes               | 3.7-     | ✅ | ✅ | ✅
| LegionMonsterVarieties           | 3.7-     | ✅ | ✅ | ✅
| LegionRanks                      | 3.7-     | ✅ | ✅ | ✅
| LegionRankTypes                  | 3.7-     | ✅ | ✅ | ✅
| LegionRewards                    | 3.7-3.11 | ❓ | ❓ | ❓
| LegionRewardTypes                | 3.7-     | ✅ | ✅ | ✅
| LegionRewardTypeVisuals          | 3.7-     | ✅ | ✅ | ✅
| LevelRelativePlayerScaling       | 3.3-     | ✅ | ✅ | ❌

## M

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| MagicMonsterLifeScalingPerLevel  | 3.7-     | ✅ | ✅ | ✅
| MapCompletionAchievements        | 3.8-     | ✅ | ✅ | ✅
| MapConnections                   | all      | ✅ | ✅ | ✅
| MapCreationInformation           | 3.5-     | ✅ | ✅ | ✅
| MapDeviceRecipes                 | 2.6-     | ❌ | ❌ | ❌
| MapDevices                       | 3.5-     | ✅ | ✅ | ❌
| MapFragmentFamilies              | 3.5-     | ✅ | ✅ | ✅
| MapFragmentMods                  | 3.2-     | ✅ | ✅ | ❌
| MapInhabitants                   | 2.6-     | ✅ | ✅ | ❌
| MapPins                          | all      | ✅ | ✅ | ❌
| MapPurchaseCosts                 | 3.5-     | ✅ | ✅ | ❌
| Maps                             | all      | ✅ | ✅ | ❌
| MapSeries                        | 3.2-     | ✅ | ✅ | ❌
| MapSeriesTiers                   | 3.6-     | ❌ | ❌ | ❌
| MapStashSpecialTypeEntries       | 3.17-    | ❓ | ❓ | ❓
| MapStashTabLayout                | 2.1-3.16 | ✅ | ✅ | ✅
| MapStatAchievements              | 3.7-3.7  | ❓ | ❓ | ❓
| MapStatConditions                | 3.8-     | ❌ | ❌ | ✅
| MapStatsFromMapStats             | 3.10-3.11 | ❓ | ❓ | ❓
| MapTierAchievements              | 3.8-     | ✅ | ✅ | ✅
| MapTiers                         | 3.7-     | ✅ | ✅ | ❌
| MasterActWeights                 | 3.0-3.4  | ❓ | ❓ | ❓
| MasterHideoutLevels              | 3.5-3.15 | ✅ | ✅ | ✅
| MavenDialog                      | 3.13-    | ❌ | ❌ | ✅
| MavenFights                      | 3.13-    | ❌ | ❌ | ❌
| Melee                            | 3.4-     | ✅ | ✅ | ❌
| MeleeTrails                      | 3.7-     | ✅ | ✅ | ❌
| MetamorphLifeScalingPerLevel     | 3.9-     | ✅ | ✅ | ✅
| MetamorphosisMetaMonsters        | 3.9-     | ✅ | ✅ | ✅
| MetamorphosisMetaSkills          | 3.9-     | ✅ | ⚠️ | ❌
| MetamorphosisMetaSkillTypes      | 3.9-     | ✅ | ✅ | ✅
| MetamorphosisRewardTypeItemsClient | 3.9-     | ✅ | ✅ | ✅
| MetamorphosisRewardTypes         | 3.9-     | ✅ | ✅ | ✅
| MetamorphosisScaling             | 3.9-     | ✅ | ✅ | ❌
| MetamorphosisStashTabLayout      | 3.11-    | ✅ | ✅ | ✅
| MicroMigrationData               | 3.0-     | ✅ | ✅ | ✅
| MicrotransactionCategory         | 3.17-    | ❓ | ❓ | ❓
| MicrotransactionCategoryId       | 3.17-    | ❓ | ❓ | ❓
| MicrotransactionCharacterPortraitVariations | 2.5-     | ✅ | ✅ | ✅
| MicrotransactionCombineFormula   | 3.10-    | ✅ | ✅ | ❌
| MicrotransactionCombineForumula  | 3.2-3.9  | ❓ | ❓ | ❓
| MicrotransactionCursorVariations | 3.13-    | ✅ | ✅ | ❌
| MicrotransactionFireworksVariations | 1.2-     | ✅ | ✅ | ✅
| MicrotransactionGemCategory      | 3.14-3.16 | ✅ | ✅ | ✅
| MicrotransactionHudVariations    | 3.13-3.13 | ❓ | ❓ | ❓
| MicrotransactionJewelVariations  | 3.17-    | ❓ | ❓ | ❓
| MicrotransactionPeriodicCharacterEffectVariations | 3.2-     | ❌ | ❌ | ❌
| MicrotransactionPortalVariations | 1.1-     | ✅ | ✅ | ❌
| MicrotransactionRarityDisplay    | 3.4-     | ✅ | ✅ | ✅
| MicrotransactionRecycleCategories | 3.8-     | ✅ | ✅ | ✅
| MicrotransactionRecycleOutcomes  | 3.8-     | ✅ | ✅ | ✅
| MicrotransactionRecycleSalvageValues | 3.8-     | ✅ | ✅ | ✅
| MicrotransactionSlot             | 3.17-    | ❓ | ❓ | ❓
| MicrotransactionSlotId           | 3.0-     | ✅ | ✅ | ✅
| MicrotransactionSocialFrameVariations | 1.1-     | ✅ | ❌ | ❌
| MinimapIcons                     | 2.1-     | ✅ | ✅ | ❌
| MiniQuestStates                  | 3.14-    | ✅ | ✅ | ✅
| MiscAnimated                     | all      | ✅ | ✅ | ❌
| MiscAnimatedArtVariations        | 3.15-    | ⚠️ | ⚠️ | ⚠️
| MiscBeams                        | 1.3-     | ✅ | ✅ | ❌
| MiscBeamsArtVariations           | 3.15-    | ⚠️ | ⚠️ | ⚠️
| MiscEffectPacks                  | 3.7-     | ✅ | ✅ | ✅
| MiscEffectPacksArtVariations     | 3.15-    | ⚠️ | ⚠️ | ⚠️
| MiscObjects                      | all      | ✅ | ✅ | ❌
| MiscObjectsArtVariations         | 3.15-    | ⚠️ | ⚠️ | ⚠️
| MissionFavourPerLevel            | 3.5-3.15 | ✅ | ✅ | ✅
| MissionTileMap                   | 1.2-     | ✅ | ✅ | ❌
| MissionTimerTypes                | 3.8-     | ✅ | ✅ | ❌
| MissionTransitionTiles           | 1.2-     | ✅ | ✅ | ✅
| ModAuraFlags                     | all      | ✅ | ✅ | ✅
| ModDomains                       | all      | ✅ | ✅ | ✅
| ModEffectStats                   | 3.12-    | ✅ | ✅ | ✅
| ModEquivalencies                 | 3.11-    | ✅ | ✅ | ❌
| ModFamily                        | all      | ✅ | ✅ | ✅
| ModGenerationType                | all      | ✅ | ✅ | ✅
| Mods                             | all      | ❌ | ❌ | ❌
| ModSellPrices                    | 0.11-3.4 | ✅ | ✅ | ✅
| ModSellPriceTypes                | 2.2-     | ✅ | ✅ | ✅
| ModSetNames                      | 3.10-    | ✅ | ✅ | ✅
| ModSets                          | 3.10-3.14 | ✅ | ✅ | ✅
| ModType                          | all      | ❌ | ❌ | ❌
| MonsterAdditionalMonsterDrops    | 2.5-3.11 | ❓ | ❓ | ❓
| MonsterArmours                   | 2.1-     | ✅ | ✅ | ✅
| MonsterBehavior                  | all      | ✅ | ✅ | ✅
| MonsterBonuses                   | 3.6-     | ✅ | ✅ | ✅
| MonsterChanceToDropItemTemplate  | 3.6-3.11 | ❓ | ❓ | ❓
| MonsterConditionalEffectPacks    | 3.12-    | ✅ | ✅ | ✅
| MonsterConditions                | 3.6-     | ✅ | ✅ | ❌
| MonsterDeathAchievements         | 3.8-     | ✅ | ✅ | ❌
| MonsterDeathConditions           | 3.8-     | ✅ | ✅ | ❌
| MonsterFleeConditions            | all      | ✅ | ✅ | ✅
| MonsterGroupEntries              | 1.3-     | ✅ | ✅ | ✅
| MonsterGroupNames                | 1.3-     | ✅ | ✅ | ✅
| MonsterHeightBrackets            | 3.10-    | ❌ | ❌ | ❌
| MonsterHeights                   | 3.10-    | ✅ | ✅ | ❌
| MonsterMapBossDifficulty         | 2.0-     | ✅ | ✅ | ❌
| MonsterMapDifficulty             | 1.3-     | ✅ | ✅ | ✅
| MonsterMortar                    | 3.2-     | ❌ | ❌ | ❌
| MonsterPackCounts                | 3.6-     | ✅ | ✅ | ✅
| MonsterPackEntries               | all      | ✅ | ✅ | ✅
| MonsterPacks                     | all      | ✅ | ✅ | ❌
| MonsterProjectileAttack          | 3.2-     | ✅ | ✅ | ❌
| MonsterProjectileSpell           | 3.2-     | ✅ | ✅ | ❌
| MonsterPushTypes                 | 3.14-    | ❓ | ❓ | ❓
| MonsterResistances               | 1.2-     | ✅ | ✅ | ✅
| MonsterScalingByLevel            | all      | ✅ | ✅ | ✅
| MonsterSegments                  | 2.1-     | ✅ | ✅ | ✅
| MonsterSize                      | all      | ✅ | ✅ | ✅
| MonsterSkillsAliveDead           | 3.4-     | ✅ | ✅ | ✅
| MonsterSkillsAttackSpell         | 3.4-     | ✅ | ✅ | ✅
| MonsterSkillsClientInstance      | 3.3-     | ✅ | ✅ | ✅
| MonsterSkillsHull                | 3.4-     | ✅ | ✅ | ✅
| MonsterSkillsOrientation         | 3.3-     | ✅ | ✅ | ✅
| MonsterSkillsPlacement           | 3.6-     | ✅ | ✅ | ✅
| MonsterSkillsReference           | 3.3-     | ✅ | ✅ | ✅
| MonsterSkillsSequenceMode        | 3.6-     | ✅ | ✅ | ✅
| MonsterSkillsShape               | 3.3-     | ✅ | ✅ | ✅
| MonsterSkillsTargets             | 3.4-     | ✅ | ✅ | ✅
| MonsterSkillsWaveDirection       | 3.8-     | ✅ | ✅ | ✅
| MonsterSpawnerGroups             | 3.5-     | ✅ | ✅ | ✅
| MonsterSpawnerGroupsPerLevel     | 3.5-     | ✅ | ✅ | ✅
| MonsterSpawnerOverrides          | 3.0-     | ✅ | ✅ | ✅
| MonsterStatsFromMapStats         | 3.6-3.11 | ❓ | ❓ | ❓
| MonsterTypes                     | all      | ✅ | ✅ | ❌
| MonsterVarieties                 | all      | ⚠️ | ⚠️ | ❌
| MonsterVarietiesArtVariations    | 3.15-    | ⚠️ | ⚠️ | ⚠️
| MouseCursorSizeSettings          | 3.15-    | ✅ | ✅ | ✅
| MoveDaemon                       | 3.4-     | ❌ | ❌ | ❌
| MTXSetBonus                      | 3.12-    | ✅ | ✅ | ❌
| MultiPartAchievementAreas        | 3.5-     | ✅ | ✅ | ❌
| MultiPartAchievementConditions   | 3.8-     | ✅ | ✅ | ✅
| MultiPartAchievements            | 3.2-     | ✅ | ✅ | ❌
| Music                            | all      | ✅ | ✅ | ❌
| MusicCategories                  | 3.6-     | ✅ | ✅ | ❌
| MysteryBoxes                     | 3.0-     | ✅ | ✅ | ❌
| MysteryPack                      | 1.3-3.1  | ❓ | ❓ | ❓
| MysteryPackItems                 | 1.3-3.1  | ❓ | ❓ | ❓

## N

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| NearbyMonsterConditions          | 3.8-     | ✅ | ✅ | ✅
| NetTiers                         | 3.2-     | ✅ | ✅ | ✅
| NormalDifficultyMasterWeights    | 1.2-2.6  | ❓ | ❓ | ❓
| Notifications                    | 1.3-     | ✅ | ✅ | ❌
| NPCAdditionalVendorItems         | 3.5-3.11 | ❓ | ❓ | ❓
| NPCAudio                         | 2.3-     | ✅ | ✅ | ✅
| NPCConversations                 | 3.12-    | ✅ | ✅ | ✅
| NPCDialogueStyles                | 3.12-    | ✅ | ✅ | ❌
| NPCFollowerVariations            | 3.5-     | ✅ | ✅ | ❌
| NPCMaster                        | 1.2-     | ✅ | ✅ | ❌
| NPCMasterExperiencePerLevel      | 1.2-3.4  | ❓ | ❓ | ❓
| NPCMasterLevels                  | 3.5-3.8  | ❓ | ❓ | ❓
| NPCPortraits                     | 3.14-    | ✅ | ✅ | ✅
| NPCs                             | all      | ✅ | ✅ | ❌
| NPCShop                          | 2.3-     | ✅ | ✅ | ❌
| NPCShopAdditionalItems           | 3.5-3.11 | ❓ | ❓ | ❓
| NPCShopSellPriceType             | 3.15-    | ❓ | ❓ | ❓
| NPCTalk                          | all      | ❌ | ❌ | ❌
| NPCTalkCategory                  | all      | ✅ | ✅ | ❌
| NPCTalkConsoleQuickActions       | 3.5-     | ✅ | ✅ | ❌
| NPCTextAudio                     | all      | ✅ | ✅ | ✅
| NPCTextAudioInterruptRules       | 3.12-    | ✅ | ✅ | ✅

## O

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| OldMapStashTabLayout             | 2.1-3.16 | ✅ | ✅ | ✅
| OnKillAchievements               | 3.8-     | ✅ | ✅ | ✅
| Orientations                     | 2.3-     | ✅ | ✅ | ✅

## P

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| PackFormation                    | 3.12-    | ✅ | ✅ | ✅
| PantheonPanelLayout              | 3.0-     | ✅ | ✅ | ❌
| PantheonSouls                    | 3.0-     | ✅ | ✅ | ❌
| PassiveJewelDistanceList         | -        | ❓ | n/a | n/a
| PassiveJewelRadii                | 3.16-    | ✅ | ✅ | ✅
| PassiveJewelSlots                | 2.0-     | ✅ | ✅ | ❌
| PassiveSkillBuffs                | 3.5-3.13 | ✅ | ✅ | ✅
| PassiveSkillFilterCatagories     | 3.11-    | ✅ | ✅ | ✅
| PassiveSkillFilterOptions        | 3.11-    | ✅ | ✅ | ✅
| PassiveSkillMasteryEffects       | 3.16-    | ✅ | ✅ | ✅
| PassiveSkillMasteryGroups        | 3.16-    | ✅ | ✅ | ✅
| PassiveSkills                    | all      | ❌ | ❌ | ❌
| PassiveSkillStatCategories       | 3.2-     | ✅ | ✅ | ✅
| PassiveSkillTrees                | 3.15-    | ❌ | ❌ | ✅
| PassiveSkillTreeTutorial         | 3.0-     | ✅ | ✅ | ❌
| PassiveSkillTreeUIArt            | 3.17-    | ❓ | ❓ | ❓
| PassiveSkillTypes                | 3.13-    | ❓ | ❓ | ❓
| PassiveTreeExpansionJewels       | 3.10-    | ✅ | ✅ | ✅
| PassiveTreeExpansionJewelSizes   | 3.10-    | ✅ | ✅ | ✅
| PassiveTreeExpansionSkills       | 3.10-    | ✅ | ✅ | ✅
| PassiveTreeExpansionSpecialSkills | 3.10-    | ✅ | ✅ | ✅
| PathOfEndurance                  | 2.6-3.11 | ❓ | ❓ | ❓
| PCBangRewardMicros               | 3.7-     | ✅ | ✅ | ✅
| PerandusBosses                   | 2.2-3.15 | ✅ | ✅ | ✅
| PerandusChests                   | 2.2-3.15 | ✅ | ✅ | ✅
| PerandusDaemons                  | 2.2-3.15 | ✅ | ✅ | ✅
| PerandusGuards                   | 2.2-3.15 | ✅ | ✅ | ❌
| PerLevelValues                   | all      | ✅ | ✅ | ✅
| Pet                              | all      | ❌ | ❌ | ✅
| PlayerConditions                 | 3.8-     | ✅ | ✅ | ❌
| PreloadGroups                    | 2.3-     | ✅ | ✅ | ✅
| PreloadPriorities                | 2.3-     | ✅ | ✅ | ✅
| PrimordialBossLifeScalingPerLevel | 3.17-    | ❓ | ❓ | ❓
| PriorityAnimatedObjects          | 3.15-3.15 | ❓ | ❓ | ❓
| ProjectileCollisionTypes         | 3.14-    | ❓ | ❓ | ❓
| Projectiles                      | all      | ⚠️ | ❌ | ❌
| ProjectilesArtVariations         | 3.15-    | ⚠️ | ⚠️ | ⚠️
| ProjectileVariations             | 2.4-3.14 | ✅ | ✅ | ✅
| Prophecies                       | 2.3-3.16 | ✅ | ✅ | ❌
| ProphecyChain                    | 2.3-3.16 | ✅ | ✅ | ❌
| ProphecySetNames                 | 3.8-3.16 | ✅ | ✅ | ✅
| ProphecySets                     | 3.8-3.11 | ❓ | ❓ | ❓
| ProphecyType                     | 2.3-3.16 | ✅ | ✅ | ❌
| PVPTypes                         | all      | ✅ | ✅ | ✅

## Q

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| Quest                            | all      | ✅ | ✅ | ❌
| QuestAchievements                | 2.1-     | ❌ | ❌ | ❌
| QuestFlags                       | all      | ✅ | ✅ | ✅
| QuestItems                       | 3.13-    | ✅ | ✅ | ✅
| QuestRewardOffers                | 3.9-     | ❌ | ❌ | ✅
| QuestRewards                     | all      | ⚠️ | ⚠️ | ❌
| QuestRewardType                  | 3.7-3.16 | ✅ | ✅ | ✅
| QuestStateCalcuation             | 3.4-3.10 | ❓ | ❓ | ❓
| QuestStateCalculation            | 3.10-    | ✅ | ✅ | ✅
| QuestStates                      | all      | ✅ | ❌ | ❌
| QuestStaticRewards               | all      | ✅ | ✅ | ❌
| QuestTrackerGroup                | 3.13-    | ✅ | ✅ | ✅
| QuestType                        | 3.5-     | ✅ | ✅ | ✅
| QuestVendorRewards               | 2.0-3.11 | ❓ | ❓ | ❓

## R

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| RaceAreas                        | 2.4-3.11 | ❓ | ❓ | ❓
| Races                            | 2.4-     | ❌ | ❌ | ❌
| RaceTimes                        | 2.6-     | ✅ | ✅ | ✅
| RandomUniqueMonsters             | 1.2-3.13 | ✅ | ❌ | ❌
| RareMonsterLifeScalingPerLevel   | 3.7-     | ✅ | ✅ | ✅
| Rarity                           | all      | ✅ | ✅ | ❌
| Realms                           | all      | ✅ | ✅ | ❌
| RecipeUnlockDisplay              | 3.5-     | ✅ | ✅ | ❌
| RecipeUnlockObjects              | 3.5-     | ✅ | ✅ | ✅
| RelativeImportanceConstants      | all      | ✅ | ✅ | ✅
| RitualBalancePerLevel            | 3.13-    | ✅ | ✅ | ✅
| RitualConstants                  | 3.13-    | ✅ | ✅ | ✅
| RitualDaemons                    | 3.13-3.13 | ❓ | ❓ | ❓
| RitualRuneTypes                  | 3.13-    | ✅ | ✅ | ❌
| RitualSetKillAchievements        | 3.13-    | ✅ | ✅ | ✅
| RitualSpawnPatterns              | 3.13-    | ✅ | ✅ | ✅
| RogueExiles                      | 2.1-     | ✅ | ✅ | ❌
| Rulesets                         | 3.11-    | ✅ | ✅ | ✅
| RunicCircles                     | 2.0-     | ✅ | ✅ | ✅

## S

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| SafehouseBYOCrafting             | 3.5-     | ✅ | ✅ | ❌
| SafehouseCraftingSpree           | 3.5-     | ✅ | ✅ | ❌
| SafehouseCraftingSpreeCurrencies | 3.5-     | ✅ | ✅ | ✅
| SalvageBoxes                     | 3.8-     | ✅ | ✅ | ✅
| Scarabs                          | 3.11-    | ✅ | ✅ | ✅
| ScarabTypes                      | 3.11-    | ✅ | ✅ | ✅
| ScoutingReports                  | 3.17-    | ❓ | ❓ | ❓
| SessionQuestFlags                | 3.4-     | ✅ | ✅ | ✅
| ShaperGuardians                  | 3.17-    | ❓ | ❓ | ❓
| ShaperMemoryFragments            | 3.2-3.8  | ❓ | ❓ | ❓
| ShaperOrbs                       | 3.2-3.8  | ❓ | ❓ | ❓
| ShieldTypes                      | all      | ✅ | ✅ | ✅
| ShopCategory                     | 0.11-3.13 | ✅ | ❌ | ❌
| ShopCountry                      | 2.1-3.13 | ✅ | ✅ | ✅
| ShopCurrency                     | 2.1-3.13 | ✅ | ✅ | ✅
| ShopForumBadge                   | 3.5-3.13 | ✅ | ✅ | ✅
| ShopItem                         | 0.11-3.1 | ❌ | ❌ | ❌
| ShopItemPrice                    | 2.0-3.1  | ❓ | ❓ | ❓
| ShopPackagePlatform              | 3.0-3.13 | ✅ | ✅ | ✅
| ShopPaymentPackage               | 0.11-3.13 | ⚠️ | ❌ | ❌
| ShopPaymentPackageItems          | 2.4-3.11 | ❓ | ❓ | ❓
| ShopPaymentPackagePrice          | 2.1-3.13 | ✅ | ✅ | ✅
| ShopPaymentPackageProxy          | 3.3-3.13 | ✅ | ✅ | ✅
| ShopRegion                       | 2.0-3.13 | ✅ | ✅ | ✅
| ShopTag                          | 3.17-    | ❓ | ❓ | ❓
| ShopToken                        | 2.4-3.13 | ✅ | ✅ | ❌
| ShrineBuffs                      | 1.0-3.16 | ✅ | ✅ | ✅
| Shrines                          | all      | ❌ | ❌ | ❌
| ShrineSounds                     | 1.0-     | ✅ | ✅ | ✅
| SigilDisplay                     | 3.5-     | ✅ | ✅ | ✅
| SingleGroundLaser                | 3.16-    | ❌ | ❌ | ✅
| SkillArtVariations               | 3.15-    | ✅ | ✅ | ✅
| SkillGemInfo                     | 3.4-     | ✅ | ✅ | ❌
| SkillGems                        | all      | ✅ | ✅ | ❌
| SkillMines                       | 3.8-     | ✅ | ✅ | ✅
| SkillMineVariations              | 3.8-     | ✅ | ✅ | ✅
| SkillMorphDisplay                | 3.7-     | ⚠️ | ❌ | ❌
| SkillMorphDisplayOverlayCondition | 3.12-    | ✅ | ✅ | ✅
| SkillMorphDisplayOverlayStyle    | 3.12-    | ✅ | ✅ | ✅
| SkillSurgeEffects                | 3.4-     | ✅ | ✅ | ❌
| SkillTotems                      | 2.1-     | ✅ | ✅ | ✅
| SkillTotemVariations             | 2.1-     | ✅ | ✅ | ✅
| SkillTrapVariations              | 3.4-     | ✅ | ✅ | ✅
| SocketNotches                    | 3.16-    | ✅ | ✅ | ✅
| SoundEffects                     | all      | ✅ | ✅ | ❌
| SpawnAdditionalChestsOrClusters  | 3.6-     | ✅ | ✅ | ✅
| SpawnObject                      | 3.3-     | ❌ | ❌ | ✅
| SpecialRooms                     | 3.6-     | ✅ | ✅ | ✅
| SpecialTiles                     | 3.6-     | ✅ | ✅ | ✅
| SpectreOverrides                 | 3.10-    | ✅ | ✅ | ✅
| StartingPassiveSkills            | 3.2-     | ✅ | ✅ | ✅
| StashId                          | 3.7-     | ✅ | ✅ | ✅
| StashTabAffinities               | 3.13-    | ✅ | ✅ | ✅
| StashType                        | 3.7-     | ✅ | ✅ | ❌
| StatDescriptionFunctions         | 1.2-     | ✅ | ✅ | ❌
| StatInterpolationTypes           | 3.0-     | ✅ | ✅ | ✅
| Stats                            | all      | ✅ | ✅ | ❌
| StatSemantics                    | all      | ✅ | ✅ | ✅
| StatSets                         | 3.10-3.11 | ✅ | ✅ | ✅
| StrDexIntMissionExtraRequirement | 1.2-3.16 | ✅ | ✅ | ✅
| StrDexIntMissionMaps             | 1.2-3.3  | ❓ | ❓ | ❓
| StrDexIntMissionMods             | 1.2-3.3  | ❓ | ❓ | ❓
| StrDexIntMissions                | 1.2-3.16 | ✅ | ✅ | ❌
| StrDexIntMissionUniqueMaps       | 1.2-3.3  | ❓ | ❓ | ❓
| StrDexMissionArchetypes          | 1.2-3.4  | ❓ | ❓ | ❓
| StrDexMissionMods                | 1.2-3.4  | ❓ | ❓ | ❓
| StrDexMissions                   | 1.2-3.4  | ❓ | ❓ | ❓
| StrDexMissionSpecialMods         | 1.2-3.4  | ❓ | ❓ | ❓
| StrIntMissionMonsterWaves        | 1.2-3.4  | ❓ | ❓ | ❓
| StrIntMissionRelicMods           | 1.2-3.4  | ❓ | ❓ | ❓
| StrIntMissionRelicPatterns       | 1.2-3.4  | ❓ | ❓ | ❓
| StrIntMissions                   | 1.2-3.4  | ❓ | ❓ | ❓
| StrMapMods                       | 1.2-2.6  | ❓ | ❓ | ❓
| StrMissionBosses                 | 1.2-3.4  | ❓ | ❓ | ❓
| StrMissionMapModNumbers          | 1.2-3.4  | ❓ | ❓ | ❓
| StrMissionMapMods                | 3.0-3.4  | ❓ | ❓ | ❓
| StrMissions                      | 1.2-3.4  | ❓ | ❓ | ❓
| StrMissionSpiritEffects          | 1.2-3.4  | ❓ | ❓ | ❓
| StrMissionSpiritSecondaryEffects | 1.2-3.4  | ❓ | ❓ | ❓
| StrongBoxEffects                 | 1.1-2.1  | ❓ | n/a | ❓
| Strongboxes                      | 2.2-     | ❌ | ❌ | ❌
| SuicideExplosion                 | 3.4-     | ✅ | ✅ | ❌
| SummonedSpecificBarrels          | 2.4-     | ❌ | ❌ | ❌
| SummonedSpecificMonsters         | 2.0-     | ❌ | ❌ | ❌
| SummonedSpecificMonstersOnDeath  | 2.4-     | ✅ | ✅ | ✅
| SuperShaperInfluence             | 3.4-3.8  | ❓ | ❓ | ❓
| SupporterPackSets                | 2.0-3.13 | ✅ | ✅ | ❌
| SurgeCategory                    | 3.9-     | ✅ | ✅ | ✅
| SurgeEffects                     | 3.14-    | ✅ | ✅ | ✅
| SurgeTypes                       | 3.4-     | ✅ | ✅ | ❌
| SynthesisAreas                   | 3.6-     | ✅ | ✅ | ✅
| SynthesisAreaSize                | 3.6-     | ✅ | ✅ | ✅
| SynthesisBonuses                 | 3.6-     | ✅ | ✅ | ✅
| SynthesisBrackets                | 3.6-     | ✅ | ✅ | ✅
| SynthesisFragmentDialogue        | 3.6-     | ✅ | ✅ | ✅
| SynthesisGlobalMods              | 3.6-     | ✅ | ✅ | ✅
| SynthesisMonsterExperiencePerLevel | 3.6-     | ✅ | ✅ | ✅
| SynthesisRewardCategories        | 3.6-     | ✅ | ✅ | ✅
| SynthesisRewardTypes             | 3.6-     | ✅ | ✅ | ✅

## T

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| TableMonsterSpawners             | 3.5-     | ✅ | ✅ | ❌
| Tags                             | all      | ✅ | ✅ | ✅
| TalismanMonsterMods              | 2.1-     | ✅ | ✅ | ✅
| TalismanPacks                    | 2.1-     | ✅ | ✅ | ❌
| Talismans                        | 2.1-     | ✅ | ✅ | ✅
| TalkingPetAudioEvents            | 3.13-    | ✅ | ✅ | ✅
| TalkingPetNPCAudio               | 3.13-    | ✅ | ✅ | ✅
| TalkingPets                      | 3.13-    | ✅ | ✅ | ✅
| TencentAutoLootPetCurrencies     | 3.8-     | ✅ | ✅ | ✅
| TencentAutoLootPetCurrenciesExcludable | 3.8-     | ✅ | ✅ | ✅
| TerrainPlugins                   | 2.4-     | ✅ | ✅ | ✅
| Tips                             | 1.2-     | ✅ | ✅ | ✅
| Topologies                       | all      | ✅ | ✅ | ✅
| TormentSpirits                   | 1.3-     | ✅ | ✅ | ✅
| TradeMarketCategory              | 3.17-    | ❓ | ❓ | ❓
| TradeMarketCategoryGroups        | 3.17-    | ❓ | ❓ | ❓
| TradeMarketCategoryListAllClass  | 3.17-    | ❓ | ❓ | ❓
| TradeMarketCategoryStyleFlag     | 3.17-    | ❓ | ❓ | ❓
| TreasureHunterMissions           | 3.3-3.16 | ✅ | ✅ | ✅
| TriggerBeam                      | 3.11-    | ✅ | ✅ | ❌
| TriggerSpawners                  | 2.4-     | ✅ | ✅ | ✅
| Tutorial                         | 3.0-     | ✅ | ✅ | ✅

## U

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| UITalkCategories                 | 3.3-     | ✅ | ✅ | ✅
| UITalkText                       | 3.3-     | ✅ | ✅ | ❌
| UltimatumEncounters              | 3.14-    | ✅ | ✅ | ✅
| UltimatumEncounterTypes          | 3.14-    | ✅ | ✅ | ✅
| UltimatumItemisedRewards         | 3.14-    | ✅ | ✅ | ✅
| UltimatumModifiers               | 3.14-    | ✅ | ✅ | ✅
| UltimatumModifierTypes           | 3.14-    | ✅ | ✅ | ✅
| UltimatumTrialMasterAudio        | 3.14-    | ✅ | ✅ | ✅
| UniqueChests                     | 1.1-     | ✅ | ✅ | ❌
| UniqueFragments                  | 3.0-3.11 | ❓ | ❓ | ❓
| UniqueJewelLimits                | 2.6-     | ✅ | ✅ | ✅
| UniqueMapInfo                    | 3.2-     | ✅ | ✅ | ✅
| UniqueMaps                       | 3.5-     | ✅ | ✅ | ❌
| UniqueSetNames                   | 1.1-     | ✅ | ✅ | ✅
| UniqueSets                       | 1.1-3.2  | ❓ | ❓ | ❓
| UniqueStashLayout                | 3.6-     | ✅ | ✅ | ✅
| UniqueStashTypes                 | 3.6-     | ❌ | ❌ | ✅

## V

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| VirtualStatContextFlags          | 3.13-    | ✅ | ✅ | ✅
| VoteState                        | all      | ✅ | ✅ | ✅
| VoteType                         | all      | ✅ | ✅ | ✅

## W

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| WarbandsGraph                    | 2.0-     | ✅ | ❌ | ✅
| WarbandsMapGraph                 | 2.0-     | ✅ | ❌ | ✅
| WarbandsPackMonsters             | 2.0-     | ✅ | ✅ | ❌
| WarbandsPackNumbers              | 2.0-     | ✅ | ✅ | ❌
| WeaponArmourCommon               | all      | ✅ | ✅ | ✅
| WeaponClasses                    | all      | ✅ | ✅ | ❌
| WeaponDamageScaling              | all      | ✅ | ✅ | ✅
| WeaponImpactSoundData            | all      | ✅ | ✅ | ✅
| WeaponSoundTypes                 | all      | ✅ | ✅ | ✅
| WeaponTypes                      | all      | ✅ | ✅ | ✅
| WindowCursors                    | 3.13-    | ✅ | ✅ | ❌
| Wordlists                        | all      | ✅ | ✅ | ✅
| Words                            | all      | ✅ | ✅ | ❌
| WorldAreaLeagueChances           | 3.16-    | ❌ | ❌ | ✅
| WorldAreas                       | all      | ✅ | ✅ | ❌
| WorldPopupIconTypes              | 3.12-    | ✅ | ✅ | ✅

## X

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| XPPerLevelForMissions            | 1.2-3.4  | ✅ | ✅ | ✅

## Z

| File                             | Releases | current  | dat64 | history
| -------------------------------- | -------- | -------- | ----- | --------
| ZanaLevels                       | 3.9-3.16 | ✅ | ✅ | ✅
| ZanaQuests                       | 3.0-3.8  | ✅ | ✅ | ✅
