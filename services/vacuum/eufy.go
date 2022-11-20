package vacuum

var (
	// home
	COMMAND_HOME = []int64{3045228, 2954657, 552182, 1448995, 549838, 476089, 523223, 477443, 469839, 1452902, 546921, 480047, 519474, 1456131, 554682, 1447641, 535619, 1455454, 555411, 471974, 527443, 472911, 526714, 473380, 525828, 474474, 524734, 1450870, 549161, 1452797, 547338, 1454777, 539577, 1452173, 548328, 478693, 520463, 480099, 519578, 480932, 528953, 471245, 527911, 472131, 527339, 472807, 526505, 473745, 525828, 479682, 519838, 1455193, 555411, 471922, 527494, 1473943, 525932, 474943, 524838, 475464, 523640, 1451495, 548797, 1453786, 524734, 481713, 528380, 30076195, 3046479, 2959865, 551348, 1450767, 549266, 477704, 521713, 478693, 520828, 1454100, 546296, 481037, 529214, 1445714, 554108, 1448475, 551556, 1455558, 555150, 472235, 527078, 473172, 526296, 474266, 525307, 474839, 524786, 1476391, 523119, 1452433, 547807, 1454568, 556088, 1451234, 548692, 478329, 520932, 479474, 519942, 480204, 519162, 481245, 528900, 471089, 528328, 471974, 527546, 472443, 526922, 478224, 521453, 1453682, 546036, 480933, 529317, 1445922, 554057, 473224, 526297, 473849, 520672, 1449203, 550828, 1451287, 548535, 479839, 519942, 29892656, 3050020, 2956375, 554994, 1447277, 552754, 474526, 525255, 475047, 524109, 1450662, 549525, 477704, 522078, 1453371, 546713, 1455141, 555359, 1451964, 548275, 478537, 520828, 479526, 520047, 479839, 519214, 480880, 518797, 1455506, 554994, 1447016, 553171, 1448735, 551140, 1455974, 555047, 471662, 527494, 472443, 526974, 473277, 525932, 473745, 525827, 474318, 524943, 475047, 524786, 475099, 524161, 480880, 529318, 1445661, 554213, 472495, 527130, 1448006, 552077, 474631, 524525, 475985, 523953, 1450766, 549057, 1452849, 547338, 480725, 529370}

	// start/stop
	COMMAND_STARTSTOP = []int64{3044340, 2955224, 552295, 1448237, 551722, 473857, 525628, 475106, 485784, 1450008, 550837, 476304, 523545, 1451935, 548857, 1453290, 522451, 1460581, 551201, 1450998, 550055, 1452456, 548128, 1454435, 546722, 1455737, 555940, 471149, 528753, 471721, 528024, 472346, 510003, 477867, 521826, 478075, 521357, 478544, 520940, 479117, 520576, 479586, 520159, 480159, 519846, 480211, 519482, 480732, 530003, 475680, 524014, 1450945, 550108, 477295, 522555, 477763, 522242, 478232, 521513, 478597, 521305, 1453967, 546982, 1455477, 538180, 1447144, 553910, 30088761, 3048454, 2958714, 553910, 1448653, 552243, 475107, 524951, 475211, 524794, 1450112, 550733, 476513, 523440, 1451987, 548649, 1453029, 547555, 1460165, 551357, 1451049, 549898, 1452352, 548493, 1453966, 546930, 1455373, 556149, 471096, 528701, 471617, 528180, 472190, 527555, 477451, 522347, 477711, 522503, 478024, 521565, 478284, 521618, 478545, 521149, 479273, 520940, 479066, 520732, 479534, 520576, 484899, 525524, 1449487, 551045, 475992, 524118, 476513, 523701, 476514, 523285, 477138, 517034, 1452039, 549118, 1453237, 547659, 1455789, 555837, 29887094, 3047725, 2959027, 553857, 1448029, 552659, 474691, 525368, 474794, 525420, 1449799, 550837, 476304, 523910, 1450945, 549639, 1452717, 548076, 1459331, 552243, 1449384, 551618, 1450633, 550107, 1452040, 549013, 1452977, 548076, 478961, 520836, 479742, 520159, 480003, 520159, 484898, 525576, 474586, 525628, 474586, 525159, 474638, 525471, 474691, 525055, 475003, 524690, 475420, 524742, 475524, 524169, 481201, 529586, 1445060, 555576, 471565, 528597, 471669, 527919, 472191, 527660, 472450, 527763, 1447248, 553701, 1448133, 552607, 1450425, 550420}

	// 30min
	COMMAND_30MIN = []int64{3038873, 2960124, 547286, 1453942, 555202, 470672, 528692, 471714, 477963, 1447484, 552494, 474474, 524942, 1450192, 549890, 1452119, 530151, 1459567, 550775, 476141, 523172, 477182, 522494, 1452432, 547338, 1454775, 555671, 471453, 528067, 1446651, 553015, 474317, 510255, 479630, 519422, 480932, 518796, 480880, 529057, 471192, 528067, 472078, 527651, 472547, 526713, 473381, 526349, 473328, 525984, 479579, 519734, 1454776, 555775, 1445661, 554161, 472755, 526765, 473484, 525932, 1448942, 554577, 446141, 549577, 1451234, 516296, 1454151, 556400, 30102899, 3049758, 2956009, 555411, 1446338, 553692, 473120, 525983, 474214, 525515, 1448942, 550775, 476037, 523693, 1451129, 548327, 1453109, 546817, 1460140, 550462, 476766, 522807, 477130, 522390, 1452223, 547494, 1454151, 556296, 470776, 528640, 1445869, 554213, 472755, 526765, 478224, 521296, 478745, 520568, 479473, 519734, 480359, 519266, 480464, 518796, 481245, 529004, 470724, 528588, 471662, 527702, 477443, 521921, 1452328, 547546, 1453734, 546452, 480567, 529474, 470619, 528588, 1446078, 549056, 472807, 526349, 1448473, 551348, 1451495, 548587, 29884568, 3044445, 2960541, 551139, 1450401, 549317, 477443, 522077, 477859, 521453, 1453109, 546556, 480047, 519473, 1454827, 555515, 1445922, 554212, 1452693, 547442, 479526, 519838, 479943, 519578, 1455140, 555462, 1445922, 554004, 472860, 526504, 1448161, 551557, 475099, 524108, 480933, 529265, 470724, 528640, 471141, 528275, 471662, 528015, 471974, 527442, 472390, 526765, 473276, 526088, 473745, 525932, 479162, 520359, 1453682, 546140, 1455349, 555462, 471610, 527494, 472547, 526869, 1447276, 552703, 474213, 525411, 1449151, 550411, 1452068, 547806}
)